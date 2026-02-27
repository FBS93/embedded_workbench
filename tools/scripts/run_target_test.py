#!/usr/bin/env python3

# ==============================================================================
# @brief Automated embedded target test runner.
#
# This script automates embedded target testing by:
# - Starting GDB server on a remote Raspberry Pi.
# - Flashing firmware via the GDB server.
# - Starting a serial-to-TCP bridge.
# - Running the target firmware.
# - Capturing serial output.
# - Reporting PASS/FAIL.
#
# Usage:
#   run_target_test.py <test_binary>
#
# Parameters:
#   <test_binary> Absolute path to the target ELF firmware to be executed.
#
# Intended for CI and automated hardware testing using CTest.
#
# @copyright
# Copyright (c) 2025 FBS93.
# See the LICENSE.md file of this project for license details.
# This notice shall be retained in all copies or substantial portions
# of the software.
#
# @warning
# This software is provided "as is", without any express or implied warranty.
# The user assumes all responsibility for its use and any consequences.
# ==============================================================================

# ==============================================================================
# IMPORTS
# ==============================================================================

# ------------------------------------------------------------------------------
# External imports
# ------------------------------------------------------------------------------
import sys
import os
import subprocess
import socket
import time

# ------------------------------------------------------------------------------
# Project-specific imports
# ------------------------------------------------------------------------------

# ==============================================================================
# CONSTANTS
# ==============================================================================

# Raspberry Pi SSH username.
RPI_USER = None

# Raspberry Pi IP address.
RPI_HOST = None

# TCP port where GDB server will listen.
GDB_PORT = None

# TCP port used for logging.
LOG_PORT = None

# Maximum allowed inactivity time (seconds) without receiving data from target.
TARGET_RX_TIMEOUT = None

# Absolute path to the workspace folder inside the container.
WORKSPACE_FOLDER = None

# Path to the target ELF firmware to be flashed and executed.
TEST_BINARY = None

# ==============================================================================
# CLASSES
# ==============================================================================

# ==============================================================================
# FUNCTIONS
# ==============================================================================

##
# @brief Retrieve a required environment variable.
#
# Reads the specified environment variable and returns its value.
# If the variable is not defined, prints an error message and
# terminates the program.
#
# @param[in] name Name of the environment variable.
# @return Value of the environment variable as a string.
##
def get_env(name):
  value = os.getenv(name)

  if value is None:
    print(f"❌ Missing required environment variable: {name}", flush=True)
    sys.exit(1)

  return value

##
# @brief Loads required configuration values.
#
# Terminates the program if any required configuration value is missing.
##
def init_config():
  global RPI_USER, RPI_HOST, GDB_PORT, LOG_PORT, TARGET_RX_TIMEOUT, WORKSPACE_FOLDER, TEST_BINARY

  # Validate command-line arguments
  if len(sys.argv) < 2:
    print("❌ Missing test binary argument.", flush=True)
    sys.exit(1)

  RPI_USER = get_env("RPI_USER")
  RPI_HOST = get_env("RPI_HOST")
  GDB_PORT = int(get_env("GDB_PORT"))
  LOG_PORT = int(get_env("LOG_PORT"))
  TARGET_RX_TIMEOUT = int(get_env("TARGET_RX_TIMEOUT"))
  WORKSPACE_FOLDER = get_env("WORKSPACE_FOLDER")
  TEST_BINARY = sys.argv[1]

##
# @brief Execute a system command.
#
# Runs the given command and optionally suppresses its output.
#
# @param[in] cmd Command to execute as a list.
# @param[in] quiet If True, suppress stdout and stderr.
##
def run(cmd, quiet=True):
  if quiet:
    subprocess.run(cmd, check=True,
                   stdout=subprocess.DEVNULL,
                   stderr=subprocess.DEVNULL)
  else:
    subprocess.run(cmd, check=True)

##
# @brief Flash firmware to the target via GDB server.
#
# Copies the test binary to the Raspberry Pi and programs it
# using GDB monitor commands.
#
# @param[in] test_binary Path to the firmware ELF file.
##
def flash_via_gdb(test_binary):
  remote_path = "/tmp/test.elf"

  print("Copying test binary to Raspberry...", flush=True)
  run(["scp", test_binary, f"{RPI_USER}@{RPI_HOST}:{remote_path}"])

  print("Flashing target via GDB server...", flush=True)
  run([
    "arm-none-eabi-gdb",
    test_binary,
    "-batch",
    "-ex", f"target remote {RPI_HOST}:{GDB_PORT}",
    "-ex", "monitor reset halt",
    "-ex", f"monitor program {remote_path} verify",
    "-ex", "quit",
  ])

  print("✅ Flash OK.", flush=True)

##
# @brief Main test execution workflow.
#
# Orchestrates the full test sequence:
# - Start GDB server
# - Flash firmware
# - Start serial bridge
# - Connect to serial output
# - Run firmware
# - Capture and analyze output
##
def main():
  init_config() # Init config values.

  print("Ensuring GDB server running...", flush=True)
  run([f"{WORKSPACE_FOLDER}/.vscode/tasks/run_target_gdb_server.sh"])

  print("Flashing firmware...", flush=True)
  flash_via_gdb(TEST_BINARY)

  print("Ensuring serial bridge running...", flush=True)
  run([f"{WORKSPACE_FOLDER}/.vscode/tasks/run_target_logging_server.sh"])

  print("Connecting to serial log stream...", flush=True)
  serial = socket.create_connection((RPI_HOST, LOG_PORT), timeout=5)

  # Flush stale serial data.
  serial.settimeout(0.2)
  try:
    while True:
      serial.recv(4096)
  except socket.timeout:
    pass

  print("Starting firmware execution...", flush=True)
  run([
    "arm-none-eabi-gdb",
    TEST_BINARY,
    "-batch",
    "-ex", f"target remote {RPI_HOST}:{GDB_PORT}",
    "-ex", "monitor reset run",
    "-ex", "quit",
  ])

  print("Capturing test output:", flush=True)
  serial.settimeout(1) # Use short timeout to periodically check for inactivity.
  buffer = ""
  last_rx = time.time() # Timestamp of last received data.

  while True:
    try:
      data = serial.recv(1024).decode(errors="ignore")
    except socket.timeout:
      data = ""

    if data:
      # Update timestamp when data is received.
      last_rx = time.time()

      print(data, end="", flush=True)
      buffer += data

      if "0 failed" in buffer:
        print("\n✅ Test pass.", flush=True)
        sys.exit(0)

      if "FAIL" in buffer:
        print("\n❌ Test fail.", flush=True)
        sys.exit(1)

    # Exit if no data has been received for TARGET_RX_TIMEOUT.
    if ((time.time() - last_rx) > TARGET_RX_TIMEOUT):
      print("\n❌ Timeout: no output from target.", flush=True)
      sys.exit(1)

# ==============================================================================
# SCRIPT ENTRY POINT
# ==============================================================================

if __name__ == "__main__":
  main()
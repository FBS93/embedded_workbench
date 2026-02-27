# Embedded target remote debugging

## Overview

This document describes how embedded target debugging is performed remotely using a Raspberry Pi acting as a remote GDB server.

The Raspberry Pi is connected to the target device through a hardware debug probe and exposes a GDB-compatible interface over TCP to the host system.

## Architecture

- Embedded target is connected to a debug probe.
- The debug probe is connected to the Raspberry Pi via USB.
- The Raspberry Pi runs a GDB server to control and debug the target.
- The GDB server is exposed over TCP.
- The host machine connects to the GDB server for remote debugging.

Target MCU --> Debug Probe --> Raspberry Pi --> GDB Server --> TCP/GDB --> Host

## Debug Setup

Remote debugging is enabled by executing  
[run_target_gdb_server.sh](../../.vscode/tasks/run_target_gdb_server.sh)

@todo Review markdown links. Hint: search all links using `](`  
@todo Document that this script can be run via VS Code task buttons.

This script:

- Connects to the Raspberry Pi via SSH.
- Ensures no stale GDB server instances are running.
- Starts the configured GDB server with the target and interface settings.
- Exposes the GDB server on the configured TCP port.
- Verifies that the GDB server is listening before returning.

Once running, the host can connect to the exposed port using GDB or compatible debug tools.

## Dependencies

The Raspberry Pi must have a compatible GDB server installed. For example:

When using OpenOCD:

1. Install OpenOCD:

    ```bash
    sudo apt-get update
    sudo apt-get install -y openocd
    ```
4. Verify installation:

    ```bash
    openocd --version
    ```

When using SEGGER J-Link GDB Server:

1. Download the latest **Linux ARM64 J-Link package** for Raspberry Pi from the official SEGGER website:

    https://www.segger.com/downloads/jlink/

2. Copy the installer to the Raspberry Pi:

    ```bash
    scp JLink_Linux_*.deb <user>@<rpi_ip>:/tmp
    ```

3. Install on the Raspberry Pi:

    ```bash
    ssh <user>@<rpi_ip>
    cd /tmp
    sudo dpkg -i JLink_Linux_*.deb
    sudo apt -f install
    ```

4. Verify installation:

    ```bash
    JLinkGDBServer -version
    ```

Make sure that the environment variables in [devcontainer.json](../../.devcontainer/devcontainer.json) are configured for the installed GDB server and target environment:

- `RPI_USER` and `RPI_HOST` must match the Raspberry Pi SSH credentials.
- `GDB_PORT` must match the port used by the GDB server.
- `SVD_FILE` must be set to the path of the SVD file used for target debugging.

Make sure that the variables in [run_target_gdb_server.sh](../../.vscode/tasks/run_target_gdb_server.sh) are configured for the installed GDB server and target environment:

- `GDB_SERVER_RUN_CMD` must be adapted to start the selected GDB server.
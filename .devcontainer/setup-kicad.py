#!/usr/bin/env python3

# ==============================================================================
# @brief Restore versioned KiCad state into the devcontainer user profile.
#
# This script restores the minimal KiCad user state stored in:
# .devcontainer/kicad-state/<version>/
#
# @copyright
# Copyright (c) 2026 FBS93.
# See the LICENSE file of this project for license details.
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
# Standard library imports
# ------------------------------------------------------------------------------
import re
import shutil
import subprocess
import sys
from pathlib import Path

# ------------------------------------------------------------------------------
# External library imports
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# Project-specific imports
# ------------------------------------------------------------------------------

# ==============================================================================
# CONSTANTS
# ==============================================================================

CONFIG_FILE_NAMES = [
  "kicad_common.json",
  "sym-lib-table",
  "fp-lib-table",
  "user.hotkeys",
  "design-block-lib-table",
]

# ==============================================================================
# FUNCTIONS
# ==============================================================================

##
# @brief Return the workspace-local KiCad state root.
#
# @return Absolute path to `.devcontainer/kicad-state`.
##
def get_state_root_dir():
  return Path(__file__).resolve().parent / "kicad-state"

##
# @brief Return the KiCad config root in the current user home.
#
# @return Absolute config root path.
##
def get_target_config_root():
  return Path.home() / ".config" / "kicad"

##
# @brief Return the KiCad share root in the current user home.
#
# @return Absolute share root path.
##
def get_target_share_root():
  return Path.home() / ".local" / "share" / "KiCad"

##
# @brief Detect the active KiCad major.minor version.
#
# @return Version string such as `10.0`, or None when detection fails.
##
def detect_kicad_version():
  config_root = get_target_config_root()
  versions = sorted(
    [path.name for path in config_root.iterdir() if path.is_dir()],
    key=parse_version_key,
  ) if config_root.exists() else []

  if versions:
    return versions[-1]

  try:
    result = subprocess.run(
      ["kicad", "--version"],
      check=True,
      capture_output=True,
      text=True,
    )
  except (OSError, subprocess.CalledProcessError):
    return None

  match = re.search(r"(\d+\.\d+)", result.stdout)
  if match is None:
    return None

  return match.group(1)

##
# @brief Convert a dotted version string into a sortable tuple.
#
# @param[in] version Version string.
# @return Tuple of integers for sorting.
##
def parse_version_key(version):
  return tuple(int(part) for part in version.split("."))

##
# @brief Remove a file or directory path when it exists.
#
# @param[in] path Path to remove.
##
def remove_path(path):
  if path.is_dir() and not path.is_symlink():
    shutil.rmtree(path)
  elif path.exists() or path.is_symlink():
    path.unlink()

##
# @brief Copy one optional state file into the KiCad profile.
#
# @param[in] source_file File in the versioned state.
# @param[in] target_file Destination file in the user profile.
##
def restore_file(source_file, target_file):
  if source_file.is_file():
    target_file.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(source_file, target_file)
  else:
    remove_path(target_file)

##
# @brief Copy one optional state directory into the KiCad profile.
#
# @param[in] source_dir Directory in the versioned state.
# @param[in] target_dir Destination directory in the user profile.
##
def restore_dir(source_dir, target_dir):
  if source_dir.is_dir():
    target_dir.parent.mkdir(parents=True, exist_ok=True)
    if target_dir.exists():
      shutil.rmtree(target_dir)
    shutil.copytree(
      source_dir,
      target_dir,
      ignore=shutil.ignore_patterns(".gitkeep"),
    )
  else:
    remove_path(target_dir)

##
# @brief Restore the versioned KiCad state when available.
##
def main():
  kicad_version = detect_kicad_version()
  if kicad_version is None:
    print("KiCad version could not be detected. Skipping KiCad state restore.")
    return 0

  state_dir = get_state_root_dir() / kicad_version
  if not state_dir.is_dir():
    print(f"No KiCad state stored for version {kicad_version}. Skipping KiCad state restore.")
    return 0

  target_config_dir = get_target_config_root() / kicad_version
  target_share_dir = get_target_share_root() / kicad_version
  state_config_dir = state_dir / "config"
  state_plugins_dir = state_dir / "data" / "plugins"
  state_template_dir = state_dir / "data" / "template"

  target_config_dir.mkdir(parents=True, exist_ok=True)
  target_share_dir.mkdir(parents=True, exist_ok=True)

  for file_name in CONFIG_FILE_NAMES:
    restore_file(state_config_dir / file_name, target_config_dir / file_name)

  restore_dir(state_config_dir / "colors", target_config_dir / "colors")
  restore_dir(state_plugins_dir, target_share_dir / "plugins")
  restore_dir(state_template_dir, target_share_dir / "template")

  print(f"KiCad state restored for version {kicad_version}.")
  return 0

# ==============================================================================
# ENTRY POINT
# ==============================================================================

if __name__ == "__main__":
  sys.exit(main())

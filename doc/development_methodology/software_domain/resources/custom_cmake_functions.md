# Custom CMake functions

## Overview

This document describes the custom CMake functions used for reusable software build integration.

## `ect_add_executable()`

[`ect_add_executable()`](../../../../tools/cmake/functions/ect_add_executable.cmake) shall be used to create executable targets that need the common executable integration defined at workspace level.

It wraps the standard CMake `add_executable()` function and applies the common executable settings configured for the workspace, avoiding per-target duplication.

The function contract and configurable variables shall be taken directly from its header documentation.

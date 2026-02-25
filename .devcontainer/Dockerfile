# Use the official Ubuntu 24.04 LTS image as the base
FROM ubuntu:24.04

# Prevent interactive prompts during package installations
ENV DEBIAN_FRONTEND=noninteractive

###############################################################################
# Common libraries and packages
###############################################################################

# Update package list and install general dependencies
#
# udev is required by VS Code Serial Monitor extension
# https://github.com/microsoft/vscode-serial-monitor/issues/137#issuecomment-1804710798
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    gdb \
    python3 \
    python3-pip \
    python3-dev \
    python-is-python3 \
    git \
    wget \
    sudo \
    locales \
    udev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Configure locale (UTF-8)
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

###############################################################################
# Install ARM GNU toolchain
###############################################################################

# ARM GNU toolchain version
ARG ARM_GCC_VERSION=15.2.rel1

# Download and install ARM GNU Toolchain
RUN wget https://developer.arm.com/-/media/Files/downloads/gnu/${ARM_GCC_VERSION}/binrel/arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz && \
    tar -xf arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz -C /opt && \
    rm arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz

# Add toolchain to PATH
ENV PATH="/opt/arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi/bin:${PATH}"

###############################################################################
# User configuration
###############################################################################

# Create non-root user, set bash as default shell, and enable passwordless sudo
RUN useradd -m -s /bin/bash fbs && \
    usermod -aG sudo fbs && \
    echo "fbs ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/fbs && \
    chmod 0440 /etc/sudoers.d/fbs

# Switch to 'fbs' user
USER fbs

# Set default workdir
WORKDIR /workspaces/embedded_c_workbench
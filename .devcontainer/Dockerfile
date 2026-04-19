FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Prevent interactive prompts during package installations
ENV DEBIAN_FRONTEND=noninteractive

################################################################################
# Common packages
################################################################################

RUN apt-get update && apt-get install -y \
    git \
    wget \
    sudo \
    locales \
    python3 \
    python3-pip \
    python3-dev \
    python-is-python3 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
    
# Install common Python tooling globally in the devcontainer so it is available
# to the whole workspace without introducing an additional virtual environment
# layer.  --break-system-packages is intentional here because this is a controlled
# development container image rather than a shared host Python installation.
RUN python -m pip install --break-system-packages \
    pytest \
    ruff

# Configure locale (UTF-8)
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

################################################################################
# Software domain
################################################################################

# udev is required by VS Code Serial Monitor extension
# https://github.com/microsoft/vscode-serial-monitor/issues/137#issuecomment-1804710798
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    clang-format \
    clang-tidy \
    gdb \
    udev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Download and install ARM GNU Toolchain
ARG ARM_GCC_VERSION=15.2.rel1
RUN wget https://developer.arm.com/-/media/Files/downloads/gnu/${ARM_GCC_VERSION}/binrel/arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz && \
    tar -xf arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz -C /opt && \
    rm arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi.tar.xz

# Add toolchain to PATH
ENV PATH="/opt/arm-gnu-toolchain-${ARM_GCC_VERSION}-x86_64-arm-none-eabi/bin:${PATH}"

# Install software-domain Python tooling globally in the devcontainer so it is 
# available to the whole workspace without introducing an additional virtual 
# environment layer.  --break-system-packages is intentional here because this 
# is a controlled development container image rather than a shared host Python 
# installation.
RUN python -m pip install --break-system-packages \
    cmakelang

################################################################################
# Hardware domain
################################################################################

# Install KiCad
RUN apt-get update && apt-get install -y \
    software-properties-common \
    && add-apt-repository --yes ppa:kicad/kicad-10.0-releases \
    && apt-get update && apt-get install -y --install-recommends \
    kicad \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install hardware-domain Python tooling globally in the devcontainer so it is
# available to the whole workspace without introducing an additional virtual
# environment layer. --break-system-packages is intentional here because this
# is a controlled development container image rather than a shared host Python
# installation.
RUN python -m pip install --break-system-packages \
    skidl \
    PySpice \
    JLC2KiCadLib

################################################################################
# Mechanical domain
################################################################################

RUN apt-get update && apt-get install -y \
    libdbus-1-3 \
    libegl1 \
    libfontconfig1 \
    libgl1 \
    libopengl0 \
    libsm6 \
    libx11-xcb1 \
    libxcb-cursor0 \
    libxcb-icccm4 \
    libxext6 \
    libxcb-keysyms1 \
    libxcb-xinerama0 \
    libxkbcommon-x11-0 \
    libxrender1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install CadQuery tooling globally in the devcontainer so it is available to the
# whole workspace without introducing an additional virtual environment layer.
# --break-system-packages is intentional here because this is a controlled
# development container image rather than a shared host Python installation.
RUN python -m pip install --break-system-packages \
    cadquery \
    cq-editor

################################################################################
# AI support with OpenCode and AI Gentle Stack
################################################################################

# Install extra dependencies required by AI Gentle Stack
RUN apt-get update && apt-get install -y \
    curl \
    nodejs \
    npm \
    golang \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Download and install AI Gentle Stack from the GitHub release tarball
ARG GENTLE_AI_VERSION=1.18.2
RUN set -eux; \
    cd /tmp; \
    curl -fsSLO "https://github.com/Gentleman-Programming/gentle-ai/releases/download/v${GENTLE_AI_VERSION}/gentle-ai_${GENTLE_AI_VERSION}_linux_amd64.tar.gz"; \
    tar -xzf "gentle-ai_${GENTLE_AI_VERSION}_linux_amd64.tar.gz"; \
    install -m 0755 gentle-ai /usr/local/bin/gentle-ai; \
    rm -f "gentle-ai_${GENTLE_AI_VERSION}_linux_amd64.tar.gz" gentle-ai

# Switch to container user so AI tooling installs in the correct user paths
USER vscode
ENV HOME=/home/vscode

# Install OpenCode 
RUN curl -fsSL https://opencode.ai/install | bash

# Execute AI Gentle Stack install command
RUN gentle-ai install \
  --preset full-gentleman \
  --agents opencode \
  --persona neutral

# Ensure Go-installed binaries are in PATH so OpenCode can launch the engram MCP server
ENV PATH="/home/vscode/.local/bin:/home/vscode/go/bin:${PATH}"

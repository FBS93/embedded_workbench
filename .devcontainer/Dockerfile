FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

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
# Install AI support with OpenCode and AI Gentle Stack
###############################################################################

# Install extra dependencies required by AI Gentle Stack
RUN apt-get update && apt-get install -y \
    curl \
    nodejs \
    npm \
    golang \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install AI Gentle Stack from the GitHub release tarball
# We avoid install.sh because it fails in Docker with "tmpdir: unbound variable" error.
ARG GENTLE_AI_VERSION=1.6.3
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
    --agent opencode \
    --preset custom \
    --component engram,sdd,skills,persona,permissions,gga \
    --skill sdd-init,sdd-explore,sdd-propose,sdd-spec,sdd-design,sdd-tasks,sdd-apply,sdd-verify,sdd-archive,skill-creator \
    --persona neutral

# Ensure Go-installed binaries are in PATH so OpenCode can launch the engram MCP server
ENV PATH="/home/vscode/.local/bin:/home/vscode/go/bin:${PATH}"
# Raspberry Pi setup

## Overview

This document describes the minimal Raspberry Pi configuration required to use it as a hardware abstraction gateway.

## Setup

For automated workflows, it is recommended to use a dedicated, locally isolated Raspberry Pi configured for this purpose with passwordless SSH access, as described below:

```bash
sudo passwd -d fbs
sudo nano /etc/ssh/sshd_config
```

Ensure the following options are set:

```text
PasswordAuthentication yes
PermitEmptyPasswords yes
```

Then restart SSH:

```bash
sudo systemctl restart ssh
```
#!/usr/bin/env bash
set -e
  
# Load and validate environment variables
: "${RPI_USER:?Missing RPI_USER}"
: "${RPI_HOST:?Missing RPI_HOST}"
: "${GDB_PORT:?Missing GDB_PORT}"

# GDB server run command.
#
# --- OpenOCD ---
# GDB_SERVER_RUN_CMD="openocd \
#   -f interface/stlink.cfg \
#   -f target/stm32f1x.cfg \
#   -c \"bindto 0.0.0.0\" \
#   -c \"gdb_port ${GDB_PORT}\""  
#
# --- SEGGER J-Link ---
GDB_SERVER_RUN_CMD="JLinkGDBServer \
  -device STM32F103C8 \
  -if SWD \
  -speed 400 \
  -port ${GDB_PORT} \
  -nogui"

ssh "$RPI_USER@$RPI_HOST" bash << EOF
set -e

# If port is already open, done.
if ss -ltn | grep -q ":$GDB_PORT"; then
    echo "✅ GDB server already listening on port $GDB_PORT."
    exit 0
fi

# Kill stale instances.
fuser -k ${GDB_PORT}/tcp 2>/dev/null || true

# Start GDB server.
nohup $GDB_SERVER_RUN_CMD > /tmp/gdb_server.log 2>&1 &

# Allow GDB server startup before port check.
sleep 2

# Wait for TCP port.
for i in {1..20}; do
    if ss -ltn | grep -q ":$GDB_PORT"; then
        echo "✅ GDB server ready on port $GDB_PORT."
        exit 0
    fi
    sleep 0.2
done

echo "❌ GDB server failed."
tail -n 20 /tmp/gdb_server.log
exit 1
EOF
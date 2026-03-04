#!/usr/bin/env bash
set -e

# -----------------------------------------------------------------------------
# This script automatically selects the FIRST device found in:
#   /dev/serial/by-id/
#
# If multiple USB/serial adapters are connected to the Raspberry Pi, this may 
# select the wrong one. In that case, manual filtering must be added.
# -----------------------------------------------------------------------------

# Load and validate environment variables
: "${RPI_USER:?Missing RPI_USER}"
: "${RPI_HOST:?Missing RPI_HOST}"
: "${LOG_PORT:?Missing LOG_PORT}"
: "${LOG_BAUD_RATE:?Missing LOG_BAUD_RATE}"

ssh "$RPI_USER@$RPI_HOST" bash << EOF
set -e

# Find serial device.
SERIAL_NAME=\$(ls -1 /dev/serial/by-id/ 2>/dev/null | head -n1)
if [ -z "\$SERIAL_NAME" ]; then
    echo "❌ No USB serial device found."
    exit 1
fi

SERIAL_DEV="/dev/serial/by-id/\$SERIAL_NAME"
echo "USB serial device found: \$SERIAL_DEV."

# Generate runtime config.
CFG=/tmp/ser2net_config.yaml

cat > \$CFG << CFGEOF
connection: &serial_dev
  accepter: tcp,$LOG_PORT
  enable: on
  options:
    banner: false
  connector: serialdev,\$SERIAL_DEV,${LOG_BAUD_RATE}n81,local
CFGEOF

# Kill previous ser2net instance.
/usr/bin/pkill ser2net 2>/dev/null || true

# Wait until the serial device node becomes available again.
for i in {1..20}; do
    if [ -e "$SERIAL_DEV" ]; then
        break
    fi
    /usr/bin/sleep 0.1
done

# Start ser2net (absolute path required in non-interactive SSH).
nohup /usr/sbin/ser2net -c \$CFG > /tmp/ser2net.log 2>&1 &

# Wait for TCP port.
for i in {1..20}; do
    if /usr/bin/ss -ltn | /usr/bin/grep -q ":$LOG_PORT"; then
        echo "✅ Serial bridge ready on TCP port $LOG_PORT."
        exit 0
    fi
    /usr/bin/sleep 0.2
done

echo "❌ ERROR: ser2net not running."
/usr/bin/tail -n 20 /tmp/ser2net.log
exit 1
EOF
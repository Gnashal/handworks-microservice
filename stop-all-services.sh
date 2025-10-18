#!/bin/bash

SERVICES=(
  "account"
  "booking"
  "inventory"
  "payment"
  # add more here when they're implemented
)

SERVICES_DIR="./services"
LOGS_DIR="./logs"
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")

mkdir -p "$LOGS_DIR"

echo "[$TIMESTAMP] Stopping selected services..." | tee -a "$LOGS_DIR/shutdown.log"

for service_name in "${SERVICES[@]}"; do
  pattern="$SERVICES_DIR/$service_name/main.go"
  pids=$(pgrep -f "$pattern")

  if [ -n "$pids" ]; then
    echo "⏹Stopping $service_name (PID(s): $pids)" | tee -a "$LOGS_DIR/shutdown.log"
    kill $pids
    sleep 1
    # Double-check if process still exists
    if pgrep -f "$pattern" > /dev/null; then
      echo "Force killing $service_name..." | tee -a "$LOGS_DIR/shutdown.log"
      pkill -9 -f "$pattern"
    fi
    echo "$service_name stopped successfully." | tee -a "$LOGS_DIR/shutdown.log"
  else
    echo "ℹ$service_name stopped." | tee -a "$LOGS_DIR/shutdown.log"
  fi
done

echo "[$TIMESTAMP] All stop operations complete." | tee -a "$LOGS_DIR/shutdown.log"

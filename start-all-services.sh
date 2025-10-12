#!/bin/bash

SERVICES=(
  "account"
  "booking"
  "inventory"
  "payment"
  # add more here when they're ready
)

SERVICES_DIR="./services"
LOGS_DIR="./logs"

# Create logs directory if not exists
mkdir -p "$LOGS_DIR"

echo "Starting selected services..."

for service_name in "${SERVICES[@]}"; do
  service_path="$SERVICES_DIR/$service_name"
  main_file="$service_path/main.go"

  if [ -f "$main_file" ]; then
    echo "Starting $service_name service..."
    nohup go run "$main_file" > "$LOGS_DIR/$service_name.log" 2>&1 &
    echo "$service_name started (PID $!) - logs: $LOGS_DIR/$service_name.log"
  else
    echo "Skipping $service_name — main.go not found."
  fi
done

echo "✅ Startup complete!"

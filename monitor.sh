#!/bin/bash

# Set thresholds (adjust as needed)
CPU_THRESHOLD=80
MEM_THRESHOLD=80

# Get CPU and memory usage
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4}')
MEM_USAGE=$(free | awk '/Mem/{printf("%.2f"), $3/$2 * 100}')

# Print system usage
echo "CPU Usage: $CPU_USAGE%"
echo "Memory Usage: $MEM_USAGE%"

# Check if usage exceeds thresholds
if (( $(echo "$CPU_USAGE > $CPU_THRESHOLD" | bc -l) )); then
    echo "Warning: High CPU Usage!"
fi

if (( $(echo "$MEM_USAGE > $MEM_THRESHOLD" | bc -l) )); then
    echo "Warning: High Memory Usage!"
fi

#!/bin/bash
echo "Starting script execution..."

# This script will perform system updates and Pi-hole maintenance automatically for my Ras-pi 5.

# 1a. Start file logging

LOG_FILE="$HOME/system_update.log"

log_message() {
    echo "$(date +"%Y-%m-%d %H:%M:%S") - $1" | sudo tee -a "$LOG_FILE"
}

# 1b. Update the Raspberry Pi OS and all installed packages.

echo "Updating package lists..."
sudo apt update
echo "Upgrading installed packages..."
sudo apt upgrade -y
echo "Cleaning up..."
sudo apt autoremove -y
sudo apt clean
echo "Update complete!"

# 2. Update Pi-hole and refresh Gravity.
#    - Use `pihole -up` to update Pi-hole to the latest version.
#    - Use `pihole -g` to update Gravity, which downloads the latest blocklists.

# 3. Log update results.
#    - Save output to a log file for debugging and record-keeping.

# 4. Schedule the script to run automatically.
#    - Add it to crontab to execute at a set interval (e.g., weekly).

#!/bin/bash
echo "Starting script execution..."

# This script will perform system updates and Pi-hole maintenance automatically for my Ras-pi 5.

# 1. Start file logging / Log update results.
#  - Save output to a log file for debugging and record-keeping.
LOG_FILE="/var/log/system_update.log"

# Echo the message to the terminal and also write it to the log file
log_message() {
    echo "$(date +"%Y-%m-%d %H:%M:%S") - $1" | sudo -E tee -a "$LOG_FILE" >/dev/null
}

# 2. Update the Raspberry Pi OS and all installed packages.
log_message "Updating package lists..."
sudo apt update | tee -a "$LOG_FILE"
log_message "Upgrading installed packages..."
sudo apt upgrade -y | tee -a "$LOG_FILE"
log_message "Cleaning up..."
sudo apt autoremove -y | tee -a "$LOG_FILE"
sudo apt clean | tee -a "$LOG_FILE"
log_message "Update complete!"

# 3. Update Pi-hole and refresh Gravity.
log_message "Updating Pi-hole..."
sudo pihole -up | tee -a "$LOG_FILE"
log_message "Refreshing Gravity..."
sudo pihole -g | tee -a "$LOG_FILE"
log_message "Pi-hole and Gravity update complete!"

#  - update Gravity, which downloads the latest blocklists.
log_message "Refreshing Gravity..."
sudo pihole -g | tee -a "$LOG_FILE"
log_message "Pi-hole and Gravity update complete!"

# 4. Schedule the script to run automatically.
#    - Add it to crontab to execute at a set interval (e.g., weekly).

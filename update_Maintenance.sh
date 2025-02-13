#!/bin/bash

# This script will perform system updates and Pi-hole maintenance automatically.

# 1. Update the Raspberry Pi OS and all installed packages.
#    - Use `sudo apt update` to fetch the latest package lists.
#    - Use `sudo apt upgrade -y` to install available updates.
#    - Optionally, use `sudo apt autoremove -y` to clean up unnecessary packages.

# 2. Update Pi-hole and refresh Gravity.
#    - Use `pihole -up` to update Pi-hole to the latest version.
#    - Use `pihole -g` to update Gravity, which downloads the latest blocklists.

# 3. Log update results.
#    - Save output to a log file for debugging and record-keeping.

# 4. Schedule the script to run automatically.
#    - Add it to crontab to execute at a set interval (e.g., weekly).

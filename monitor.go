package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// This script will monitor system metrics such as CPU temperature, memory usage, and disk utilization.

// 1. Retrieve CPU Temperature.
//    - Use `vcgencmd measure_temp` to get the Raspberry Pi’s CPU temperature.

func getCPUTemperature() (string, error) {
	// Execute the command to get the CPU temperature
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Convert output to a string and clean it up
	temp := strings.TrimSpace(string(output))
	return temp, nil
}

// 2. Monitor CPU usage.
//    - Extract CPU usage from the `top` command.

// Function to get CPU usage
func getCPUUsage() (string, error) {
	cmd := exec.Command("top", "-bn1") // Run `top` in batch mode for one iteration
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Cpu(s):") {
			return strings.TrimSpace(line), nil
		}
	}

	return "CPU usage not found", nil
}

// 3. Monitor Memory usage.
//    - Use `free -m` to check memory availability and calculate usage percentage.

// 4. Monitor Disk Utilization.
//    - Use `df -h` to check available disk space.

// 5. Display results in a formatted output.
//    - Print the values in a human-readable format.

// 6. Set up warnings if thresholds are exceeded.
//    - Alert if CPU temp is too high.
//    - Alert if memory usage exceeds 80%.
//    - Alert if disk space is critically low.

// 7. Optionally, log the results.
//    - Save monitoring data to a file for tracking system performance over time.

// 8. Schedule automatic execution.
//    - Run this script on a schedule using a cron job.

func main() {
	// Retrieve and display CPU temperature
	temp, err := getCPUTemperature()
	if err != nil {
		fmt.Println("Error retrieving CPU temperature:", err)
		return
	}

	fmt.Println("CPU Temperature:", temp)

	// Get CPU usage
	cpuUsage, err := getCPUUsage()
	if err != nil {
		fmt.Println("Error retrieving CPU usage:", err)
	} else {
		fmt.Println("CPU Usage:", cpuUsage)
	}
}

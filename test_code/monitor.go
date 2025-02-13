package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// Get CPU usage
	cpuCmd := "top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'"
	cpuUsage, err := exec.Command("bash", "-c", cpuCmd).Output()
	if err != nil {
		fmt.Println("Error fetching CPU usage:", err)
	} else {
		cpuStr := strings.TrimSpace(string(cpuUsage))
		fmt.Println("CPU Usage Command Output:", cpuStr) // Debug print
		cpuFloat, _ := strconv.ParseFloat(cpuStr, 64)

		// Get memory usage
		memCmd := "free | awk '/Mem/{printf(\"%.2f\"), $3/$2 * 100}'"
		memUsage, err := exec.Command("bash", "-c", memCmd).Output()
		if err != nil {
			fmt.Println("Error fetching memory usage:", err)
		} else {
			memStr := strings.TrimSpace(string(memUsage))
			fmt.Println("Memory Usage Command Output:", memStr) // Debug print
			memFloat, _ := strconv.ParseFloat(memStr, 64)

			// Get disk usage (use df to get free space on root partition)
			diskCmd := "df -h | awk '$NF==\"/\"{print $5}'"
			diskUsage, err := exec.Command("bash", "-c", diskCmd).Output()
			if err != nil {
				fmt.Println("Error fetching disk usage:", err)
			} else {
				diskStr := strings.TrimSpace(string(diskUsage))
				fmt.Println("Disk Usage Command Output:", diskStr) // Debug print

				// Get CPU temperature (specific for Raspberry Pi)
				tempCmd := "vcgencmd measure_temp"
				tempOutput, err := exec.Command("bash", "-c", tempCmd).Output()
				if err != nil {
					fmt.Println("Error fetching CPU temperature:", err)
				} else {
					// Debug print for temperature
					fmt.Println("Temperature Command Output:", string(tempOutput))
					// Remove the "temp=" part and extract the numeric value
					tempStr := strings.TrimSpace(string(tempOutput))
					tempStr = strings.Replace(tempStr, "temp=", "", 1) // Removes the "temp=" prefix
					tempStr = strings.Replace(tempStr, "'C", " C", 1)
					fmt.Printf("Current CPU Temperature: %s\n", tempStr)
				}

				// Print results
				fmt.Printf("CPU Usage: %.2f%%\n", cpuFloat)
				fmt.Printf("Memory Usage: %.2f%%\n", memFloat)
				fmt.Printf("Disk Usage: %s\n", diskStr)

				// Check thresholds
				if cpuFloat > 80 {
					fmt.Println("Warning: High CPU Usage!")
				}
				if memFloat > 80 {
					fmt.Println("Warning: High Memory Usage!")
				}
			}
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

// Function to calculate CPU usage over multiple samples
func getAverageCPUUsage(samples int, interval time.Duration) (string, error) {
	var totalUsage float64

	for i := 0; i < samples; i++ {
		idle1, total1 := readCPUSample()
		time.Sleep(interval) // Wait between samples
		idle2, total2 := readCPUSample()

		idleDelta := idle2 - idle1
		totalDelta := total2 - total1

		if totalDelta == 0 {
			continue
		}

		// Calculate CPU usage for this sample
		usage := 100.0 * (1.0 - float64(idleDelta)/float64(totalDelta))
		totalUsage += usage
	}

	// Compute the average CPU usage
	averageUsage := totalUsage / float64(samples)
	return fmt.Sprintf("Average CPU Usage (over %d samples): %.2f%%", samples, averageUsage), nil
}

// Reads CPU stats from /proc/stat
func readCPUSample() (idle, total int64) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if fields[0] == "cpu" {
			var values []int64
			for _, field := range fields[1:] {
				val, _ := strconv.ParseInt(field, 10, 64)
				values = append(values, val)
			}
			// Total CPU time is the sum of all fields
			total = sum(values)
			// Idle time is the 4th field
			idle = values[3]
			break
		}
	}

	return idle, total
}

// Helper function to sum integer slices
func sum(nums []int64) int64 {
	var total int64
	for _, num := range nums {
		total += num
	}
	return total
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
	cpuUsage, err := getAverageCPUUsage(5, 1*time.Second) // 5 samples over 5 seconds
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(cpuUsage)
}

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

// Thresholds
const (
	CPU_TEMP_THRESHOLD     = 45.0 // 80°C
	MEMORY_USAGE_THRESHOLD = 9.0  // 80%
	DISK_USAGE_THRESHOLD   = 10.0 // 80%
)

// 1. Retrieve CPU Temperature.
//    - Use `vcgencmd measure_temp` to get the Raspberry Pi’s CPU temperature.

func getCPUTemperature() (float64, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Example output: "temp=45.6'C"
	tempStr := strings.TrimSpace(string(output))
	tempStr = strings.TrimPrefix(tempStr, "temp=") // Remove "temp="
	tempStr = strings.TrimSuffix(tempStr, "'C")    // Remove "'C"

	// Convert to float
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	return temp, nil
}

// 2. Monitor CPU usage.
//    - Extract CPU usage using mpstat

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

func getMemoryUsage() (string, error) {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "Memory usage not found", nil
	}

	fields := strings.Fields(lines[1]) // Extract memory values from the second line
	if len(fields) < 3 {
		return "Memory usage not found", nil
	}

	totalMem, _ := strconv.Atoi(fields[1])
	usedMem, _ := strconv.Atoi(fields[2])
	memUsagePercent := (float64(usedMem) / float64(totalMem)) * 100

	return fmt.Sprintf("Memory Usage: %dMB / %dMB (%.2f%%)", usedMem, totalMem, memUsagePercent), nil
}

// 4. Monitor Disk Utilization.
//    - Use `df -h` to check available disk space.

func getDiskUsage() (string, error) {
	// Run the `df -h /` command to get disk usage for the root partition
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output format")
	}

	// Extract values from the second line (root partition info)
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return "", fmt.Errorf("could not parse disk usage")
	}

	// Extract total, used, available, and usage percentage
	total := fields[1]     // Total disk size
	used := fields[2]      // Used space
	available := fields[3] // Available space
	usage := fields[4]     // Percentage used

	// Return formatted disk usage information
	return fmt.Sprintf("Disk Usage: %s used / %s total (%s available, %s used)", used, total, available, usage), nil
}

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
	// CPU Temperature Check
	temp, err := getCPUTemperature()
	if err != nil {
		fmt.Println("Error retrieving CPU temperature:", err)
	} else {
		fmt.Printf("CPU Temperature: %.2f°C\n", temp)
		if temp >= CPU_TEMP_THRESHOLD {
			fmt.Println("WARNING: CPU temperature is too high!")
		}
	}

	// CPU Usage Check
	cpuUsage, err := getAverageCPUUsage(5, 1*time.Second)
	if err != nil {
		fmt.Println("Error retrieving CPU usage:", err)
	} else {
		fmt.Println(cpuUsage)
	}

	// Memory Usage Check
	memUsage, err := getMemoryUsage()
	if err != nil {
		fmt.Println("Error retrieving Memory usage:", err)
	} else {
		fmt.Println(memUsage)

		// Extract percentage value from memory usage string
		memParts := strings.Fields(memUsage)
		if len(memParts) > 3 {
			percentStr := strings.TrimSuffix(memParts[3], "%")
			memPercent, err := strconv.ParseFloat(percentStr, 64)
			if err == nil && memPercent >= MEMORY_USAGE_THRESHOLD {
				fmt.Println("WARNING: Memory usage is too high!")
			}
		}
	}

	// Disk Usage Check
	diskUsage, err := getDiskUsage()
	if err != nil {
		fmt.Println("Error retrieving Disk usage:", err)
	} else {
		fmt.Println(diskUsage)

		// Extract percentage value from disk usage string
		diskParts := strings.Fields(diskUsage)
		if len(diskParts) > 3 {
			percentStr := strings.TrimSuffix(diskParts[3], "%")
			diskPercent, err := strconv.ParseFloat(percentStr, 64)
			if err == nil && diskPercent >= DISK_USAGE_THRESHOLD {
				fmt.Println("WARNING: Disk usage is too high!")
			}
		}
	}
}

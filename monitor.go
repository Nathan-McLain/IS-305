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

func getMemoryUsage() (float64, int, int, error) {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("memory usage not found")
	}

	fields := strings.Fields(lines[1]) // Extract memory values from the second line
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("memory usage not found")
	}

	totalMem, _ := strconv.Atoi(fields[1])
	usedMem, _ := strconv.Atoi(fields[2])
	memUsagePercent := (float64(usedMem) / float64(totalMem)) * 100

	return memUsagePercent, usedMem, totalMem, nil
}

// 4. Monitor Disk Utilization.
//    - Use `df -k` to check available disk space in KB.

func getDiskUsage() (float64, float64, float64, error) {
	cmd := exec.Command("df", "-k", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("unexpected output format")
	}

	fields := strings.Fields(lines[1])
	totalDisk, _ := strconv.ParseFloat(fields[1], 64)
	usedDisk, _ := strconv.ParseFloat(fields[2], 64)

	diskPercent := (usedDisk / totalDisk) * 100
	return diskPercent, usedDisk / 1024 / 1024, totalDisk / 1024 / 1024, nil
}

// 6. Set up warnings if thresholds are exceeded.
//    - Alert if CPU temp is too high.
//    - Alert if memory usage exceeds 80%.
//    - Alert if disk space is critically low.

// 7. Optionally, log the results.
//    - Save monitoring data to a file for tracking system performance over time.

// 8. Schedule automatic execution.
//    - Run this script on a schedule using a cron job.

func main() {
	temp, err := getCPUTemperature()
	if err != nil {
		fmt.Println("Error retrieving CPU temperature:", err)
	} else {
		fmt.Printf("CPU Temperature: %.2f°C\n", temp)
		if temp >= CPU_TEMP_THRESHOLD {
			fmt.Println("WARNING: CPU temperature is too high!")
		}
	}

	cpuUsage, err := getAverageCPUUsage(5, 1*time.Second)
	if err != nil {
		fmt.Println("Error retrieving CPU usage:", err)
	} else {
		fmt.Println(cpuUsage)
	}

	memUsage, usedMem, totalMem, err := getMemoryUsage()
	if err != nil {
		fmt.Println("Error retrieving Memory usage:", err)
	} else {
		fmt.Printf("Memory Usage: %.2f%% (%dMB / %dMB)\n", memUsage, usedMem, totalMem)
		if memUsage >= MEMORY_USAGE_THRESHOLD {
			fmt.Println("WARNING: Memory usage is too high!")
		}
	}

	diskUsage, usedDisk, totalDisk, err := getDiskUsage()
	if err != nil {
		fmt.Println("Error retrieving Disk usage:", err)
	} else {
		fmt.Printf("Disk Usage: %.2f%% (%.2fGB / %.2fGB)\n", diskUsage, usedDisk, totalDisk)
		if diskUsage >= DISK_USAGE_THRESHOLD {
			fmt.Println("WARNING: Disk usage is too high!")
		}
	}
}

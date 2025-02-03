package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func test() {
	fmt.Println("Hello, Go!")
}

func main() {
	test()
	// Get CPU usage
	cpuCmd := "top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'"
	cpuUsage, _ := exec.Command("bash", "-c", cpuCmd).Output()
	cpuStr := strings.TrimSpace(string(cpuUsage))
	cpuFloat, _ := strconv.ParseFloat(cpuStr, 64)

	// Get memory usage
	memCmd := "free | awk '/Mem/{printf(\"%.2f\"), $3/$2 * 100}'"
	memUsage, _ := exec.Command("bash", "-c", memCmd).Output()
	memStr := strings.TrimSpace(string(memUsage))
	memFloat, _ := strconv.ParseFloat(memStr, 64)

	// Print results
	fmt.Printf("CPU Usage: %.2f%%\n", cpuFloat)
	fmt.Printf("Memory Usage: %.2f%%\n", memFloat)

	// Check thresholds
	if cpuFloat > 80 {
		fmt.Println("Warning: High CPU Usage!")
	}
	if memFloat > 80 {
		fmt.Println("Warning: High Memory Usage!")
	}
}

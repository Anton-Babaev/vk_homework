package main

import (
	"fmt"
	"strings"
)

const (
	secondsInMinute = 60
	minutesInHour   = 60
	hoursInDay      = 24
	secondsInHour   = minutesInHour * secondsInMinute
	secondsInDay    = hoursInDay * secondsInHour
)

func printUptime(total int) {
	if total == 0 {
		fmt.Println("0с")
		return
	}

	days := total / secondsInDay
	hours := (total % secondsInDay) / secondsInHour
	minutes := (total % secondsInHour) / secondsInMinute
	seconds := total % secondsInMinute

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dд", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dч", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dм", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%dс", seconds))
	}
	fmt.Println(strings.Join(parts, " "))
}

func main() {
	printUptime(93784)
	printUptime(86405)
	printUptime(3600)
	printUptime(59)
	printUptime(0)
}

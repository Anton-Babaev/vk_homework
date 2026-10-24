package main

import "fmt"

func main() {
	total := 93784

	if total == 0 {
		fmt.Println("0с")
		return
	}

	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60
	seconds := total % 60

	if days > 0 {
		fmt.Printf("%dд ", days)
	}
	if hours > 0 {
		fmt.Printf("%dч ", hours)
	}
	if minutes > 0 {
		fmt.Printf("%dм ", minutes)
	}
	if seconds > 0 {
		fmt.Printf("%dс", seconds)
	}
	fmt.Println()
}

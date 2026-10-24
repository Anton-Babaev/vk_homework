package main

import "fmt"

func main() {
	total := 93784

	// Особый случай: ноль. Без этой проверки все if ниже ложны,
	// и выведется пустая строка вместо "0с".
	if total == 0 {
		fmt.Println("0с")
		return
	}

	// Целочисленное деление и остаток.
	// 86400 — секунд в сутках, 3600 — в часе, 60 — в минуте.
	days := total / 86400           // целых дней
	hours := (total % 86400) / 3600 // остаток -> часы
	minutes := (total % 3600) / 60  // остаток -> минуты
	seconds := total % 60           // остаток -> секунды

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

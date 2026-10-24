package main

import "fmt"

const (
	Initial    = 100000.0
	AnnualRate = 0.12
)

func main() {
	balance := Initial
	monthlyRate := AnnualRate / 12

	for month := 1; month <= 12; month++ {
		balance *= (1 + monthlyRate)
		fmt.Printf("месяц %2d: %.2f\n", month, balance)
	}

	balance = Initial
	months := 0
	for balance < Initial*2 {
		balance *= (1 + monthlyRate)
		months++
	}
	fmt.Printf("Вклад удвоится через %d мес. Баланс: %.2f\n", months, balance)
}

// Почему деньги не хранят во float64:
// float64 — двоичная плавающая точка. Числа вроде 0.1
// не представимы точно, ошибка накапливается при множестве операций.
// В реальных системах используют int64 (копейки) или decimal-библиотеки,
// а в БД — NUMERIC/DECIMAL, не FLOAT.

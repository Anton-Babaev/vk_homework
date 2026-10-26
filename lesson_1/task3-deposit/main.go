package main

import "fmt"

const (
	initialDeposit = 100000.0
	annualRate     = 0.12
	periodsPerYear = 12 // Периоды капитализации
	periodToShow   = 12 // Количество показываемых периодов
)

func main() {
	balance := initialDeposit
	monthlyRate := annualRate / periodsPerYear

	for month := 1; month <= periodToShow; month++ {
		balance *= (1 + monthlyRate)
		fmt.Printf("месяц %2d: %.2f\n", month, balance)
	}

	balance = initialDeposit
	months := 0
	for balance < initialDeposit*2 {
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

package main

import (
	"errors"
	"fmt"
)

// NewAccount создаёт счёт с начальным балансом (в копейках).
// Возвращает три замыкания над общим current.
func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	current := initial

	deposit = func(amount int64) {
		current += amount
	}

	withdraw = func(amount int64) error {
		if amount > current {
			return errors.New("недостаточно средств")
		}
		current -= amount
		return nil
	}

	balance = func() int64 {
		return current
	}

	return deposit, withdraw, balance
}

// formatRub переводит копейки в рубли с двумя знаками: 150050 -> "1500.50".
func formatRub(coins int64) string {
	rubles := coins / 100
	cents := coins % 100
	return fmt.Sprintf("%d.%02d", rubles, cents)
}

func main() {
	depositA, withdrawA, balanceA := NewAccount(100000)
	_, _, balanceB := NewAccount(0)

	fmt.Println("Счёт A:", formatRub(balanceA()))

	depositA(50000)
	fmt.Println("Пополнение 500.00 ->", formatRub(balanceA()))

	if err := withdrawA(200000); err != nil {
		fmt.Printf("Снятие 2000.00: %v (баланс %s)\n", err, formatRub(balanceA()))
	}

	if err := withdrawA(30050); err != nil {
		fmt.Println("Снятие 300.50:", err)
	} else {
		fmt.Println("Снятие 300.50 ->", formatRub(balanceA()))
	}

	fmt.Println("Счёт B:", formatRub(balanceB()))
}

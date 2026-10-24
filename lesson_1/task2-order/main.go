package main

import "fmt"

// CustomerStatus — свой тип на базе int.
// Не алиас (= int), а новый именованный тип: нельзя случайно
// передать обычный int туда, где ожидается CustomerStatus.
type CustomerStatus int

// Константы статусов через iota: Regular = 0, Silver = 1, Gold = 2.
// iota — счётчик строк в блоке const, начинается с 0.
const (
	Regular CustomerStatus = iota
	Silver
	Gold
)

// String реализует интерфейс fmt.Stringer неявно: достаточно
// метода String() string — и fmt будет печатать "Gold", а не 2.
func (s CustomerStatus) String() string {
	switch s {
	case Regular:
		return "Regular"
	case Silver:
		return "Silver"
	case Gold:
		return "Gold"
	default:
		return "Unknown"
	}
}

const (
	DeliveryPrice   = 299.0
	FreeDeliveryMin = 3000.0
)

func main() {
	amount := 4500.0
	status := Gold

	var discount float64
	switch status {
	case Regular:
		discount = 0
	case Silver:
		discount = amount * 0.05
	case Gold:
		discount = amount * 0.10
	default:
		fmt.Println("неизвестный статус")
		discount = 0
	}

	afterDiscount := amount - discount

	delivery := DeliveryPrice
	if afterDiscount >= FreeDeliveryMin {
		delivery = 0
	}

	total := afterDiscount + delivery

	fmt.Printf("Сумма: %.2f, статус: %v\n", amount, status)
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f\n", total)
}

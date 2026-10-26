package main

import "fmt"

type CustomerStatus int

const (
	Regular CustomerStatus = iota
	Silver
	Gold
)

func (s CustomerStatus) String() string {
	switch s {
	case Regular:
		return "Regular"
	case Silver:
		return "Silver"
	case Gold:
		return "Gold"
	default:
		return "Неизвестный статус"
	}
}

const (
	deliveryPrice   = 299.0
	freeDeliveryMin = 3000.0
)

func printOrder(amount float64, status CustomerStatus) {
	var discount float64
	switch status {
	case Regular:
		discount = 0
	case Silver:
		discount = amount * 0.05
	case Gold:
		discount = amount * 0.10
	default:
		discount = 0
	}

	afterDiscount := amount - discount

	delivery := deliveryPrice
	if afterDiscount >= freeDeliveryMin {
		delivery = 0
	}

	total := afterDiscount + delivery

	fmt.Printf("Сумма: %.2f, статус: %v\n", amount, status)
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f\n", total)
}

func main() {
	printOrder(4500, Gold)
	fmt.Println()
	printOrder(2900, Silver)
	fmt.Println()
	printOrder(1000, Regular)
	fmt.Println()
}

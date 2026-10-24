package main

import "fmt"

type CustomerStatus int

const (
	Regular CustomerStatus = iota
	Silver
	Gold
)

const (
	DeliveryPrice   = 299.0
	FreeDeliveryMin = 3000.0
	SilverDiscount  = 0.05
	GoldDiscount    = 0.10
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
		return "Неизвестнй статус"
	}
}

func main() {
	amount := 2500.0
	status := Gold

	var discount float64
	switch status {
	case Regular:
		discount = 0
	case Silver:
		discount = amount * SilverDiscount
	case Gold:
		discount = amount * GoldDiscount
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

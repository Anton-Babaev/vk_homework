package main

import (
	"fmt"
	"time"
)

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}
	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
}

func main() {
	fmt.Println("--- n = 3 ---")
	step := 0
	Hanoi(3, "A", "C", "B", func(from, to string) {
		step++
		fmt.Printf("%d: %s -> %s\n", step, from, to)
	})

	fmt.Println("--- n = 20 ---")
	start := time.Now()
	defer func() {
		fmt.Printf("время: %s\n", time.Since(start))
	}()

	moves := 0
	Hanoi(20, "A", "C", "B", func(from, to string) {
		moves++
	})
	fmt.Printf("n=20: %d ходов\n", moves)

	fmt.Println("Вопрос: почему ходов ровно 2^n - 1")
	fmt.Println("Вопрос, почсему использовали знак ^, я думал тут про степень )))")
	fmt.Println("Ответ, потому что тут прослеживатся зависимость количества дисков к количеству")
	fmt.Println("ходов. А именно количество ходов всегда в 2 раза больше чем дисков и минус 1. ")
	fmt.Println("Минус 1, потому что каждый диск прекидываем по 2 раза, а опследний один")
}

package main

import "fmt"

const N = 8

func main() {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if (i+j)%2 == 0 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

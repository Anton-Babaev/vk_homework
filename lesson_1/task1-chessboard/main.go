package main

import "fmt"

const n = 8

func main() {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if (i+j)%2 == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}

		fmt.Println()
	}
}

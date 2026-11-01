package main

import "fmt"

func FirstUnique(s string) rune {
	count := make(map[rune]int)
	for _, r := range s {
		count[r]++
	}
	for _, r := range s {
		if count[r] == 1 {
			return r
		}
	}
	return 0
}

func main() {
	fmt.Printf("%q\n", FirstUnique("abacabad"))
	fmt.Printf("%q\n", FirstUnique("abcabcd"))
	fmt.Printf("%q\n", FirstUnique("абракадабра"))
	fmt.Printf("%q\n", FirstUnique("aabbcc"))
}

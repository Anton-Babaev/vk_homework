package main

import (
	"fmt"
	"unicode"
)

func AreAnagrams(s1, s2 string) bool {
	count := make(map[rune]int)

	for _, r := range s1 {
		if unicode.IsSpace(r) {
			continue
		}
		count[unicode.ToLower(r)]++
	}
	for _, r := range s2 {
		if unicode.IsSpace(r) {
			continue
		}
		count[unicode.ToLower(r)]--
	}

	for _, v := range count {
		if v != 0 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("listen - silent =", AreAnagrams("listen", "silent"))
	fmt.Println("Eleven plus two - Twelve plus one =", AreAnagrams("Eleven plus two", "Twelve plus one"))
	fmt.Println("Апельсин - спаниель =", AreAnagrams("Апельсин", "спаниель"))
	fmt.Println("Hello - World =", AreAnagrams("Hello", "World"))
}

package main

import (
	"fmt"
	"unicode"
)

func IsPalindrome(s string) bool {
	var cleaned []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
		if cleaned[i] != cleaned[j] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(IsPalindrome("шалаш"))
	fmt.Println(IsPalindrome("А роза упала на лапу Азора"))
	fmt.Println(IsPalindrome("Was it a car or a cat I saw?"))
	fmt.Println(IsPalindrome("hello"))
}

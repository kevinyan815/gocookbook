package main

import (
	"fmt"
	"unicode"
)

func verifyPasswordComplexity(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func main() {
	fmt.Println(verifyPasswordComplexity("pass"))     // false
	fmt.Println(verifyPasswordComplexity("password")) // false
	fmt.Println(verifyPasswordComplexity("Password")) // false
	fmt.Println(verifyPasswordComplexity("P@ssword")) // false
	fmt.Println(verifyPasswordComplexity("P@ssw0rd")) // true
}

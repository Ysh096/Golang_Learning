package main

import "fmt"

func canIDrink(age int) bool {
	switch koreanAge := age + 2; {
	case koreanAge <= 18:
		return false
	case koreanAge > 18:
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(17))
}

package main

import "fmt"

// struct 정의
type person struct {
	name    string
	age     int
	favFood []string
}

func main() {
	// value 할당 방법1
	// favFood := []string{"kimchi", "ramen"}
	// nico := person{"nico", 18, favFood}
	// fmt.Println(nico.name)

	// value 할당 방법2
	favFood := []string{"kimchi", "ramen"}
	nico := person{name: "nico", age: 18, favFood: favFood}
	fmt.Println(nico.name)
}

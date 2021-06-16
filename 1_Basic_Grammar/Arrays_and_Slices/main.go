package main

import "fmt"

func main() {
	names := []string{"nico", "lynn", "dal"}
	fmt.Println(names)
	names = append(names, "flynn")
	fmt.Println(names)
}

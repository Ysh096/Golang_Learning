package main

import "fmt"

func main() {
	nico := map[string]string{"name": "nico", "age": "12"} //map의 key는 string, value도 string
	fmt.Println(nico)
	for key, value := range nico {
		fmt.Println(key, value)
	}
}

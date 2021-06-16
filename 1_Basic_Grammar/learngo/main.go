package main

import (
	"fmt"
	"strings"
)

// func lenAndUpper(name string) (int, string) {
// 	return len(name), strings.ToUpper(name)
// }

func lenAndUpper(name string) (length int, uppercase string) {
	defer fmt.Println("I'm done") //function이 끝나고 나면 실행됨
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}

func main() {
	totalLength, up := lenAndUpper("nico")
	fmt.Println(totalLength, up)
}

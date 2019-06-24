package main

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d %d\n", width, height)
}

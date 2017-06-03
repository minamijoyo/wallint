package main

import (
	"fmt"
	"strings"

	"github.com/Pallinder/go-randomdata"
)

func main() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 5; j++ {
			str := fmt.Sprintf("%s := %d; ", strings.ToLower(randomdata.SillyName()), i*10+j)
			fmt.Print(str)
		}
		fmt.Println("")
	}
}

package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	var str string
	flag.StringVar(&str, "s", "", "String after -s")
	flag.Parse()
	words := splitString(str)
	normalized, err := stemming(words)
	if err != nil {
		fmt.Println("Something went wrong")
	}
	fmt.Println(strings.Join(normalized, " "))
}

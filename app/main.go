package main

import (
	"fmt"
)

func main() {
	var name string
	
	for true {
		fmt.Print("$ ")
		fmt.Scanf("%s\n", &name)
		fmt.Printf("%s: command not found\n", name)
	}
	
}

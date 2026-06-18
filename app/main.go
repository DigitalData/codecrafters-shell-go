package main

import (
	"fmt"
)

const CMD_EXIT = "exit"

func main() {
	var name string
	
	for true {
		fmt.Print("$ ")
		fmt.Scanf("%s\n", &name)
		if CMD_EXIT == name {
			break
		}
		fmt.Printf("%s: command not found\n", name)
	}
	
}

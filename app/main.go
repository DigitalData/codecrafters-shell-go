package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

const DEBUG = false
const CMD_EXIT = "exit"
const CMD_ECHO = "echo"

func main() {

	for true {
		
		fmt.Print("$ ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		raw_args := scanner.Text()
	
		if len(raw_args) == 0 { continue }
		cmd, cmd_args, has_args := strings.Cut(raw_args, " ")

		if DEBUG {
			fmt.Printf("DEBUG / cmd = \"%s\"\n", cmd)
			fmt.Printf("DEBUG / cmd_args = \"%s\"\n", cmd_args)
			fmt.Printf("DEBUG / has_args = %v\n", has_args)
		}

		if cmd == CMD_EXIT { break }

		switch cmd {
			case CMD_ECHO: 
				fmt.Printf("%s\n", cmd_args)
			default:
				fmt.Printf("%s: command not found\n", cmd)
		}

	}
	
}

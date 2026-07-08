package main

import (
	"bufio"
	"fmt"
	"os"
)

type CommandRune int
const (
	KeyCtrlC 		CommandRune = 3
	KeyTab 			CommandRune = 9
	KeyCtrlJ 		CommandRune = 10
	KeyEnter 		CommandRune = 13
	KeyBackspace 	CommandRune = 127
)

func read_line() string {
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)
	double_tab := false
	var line string
	for {
		next_rune, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if next_rune != rune(KeyTab) {
			double_tab = false
		}

		switch next_rune {
		case rune(KeyTab):
			line, double_tab = handle_autocomplete(line, double_tab)
		case rune(KeyBackspace):
			line_len := len(line)
			if line_len > 0 {
				fmt.Print("\b \b")
				line = line[:line_len-1]
			}
		case rune(KeyCtrlC):
			os.Exit(0)
		case rune(KeyCtrlJ), rune(KeyEnter), rune(0):
			fmt.Print("\r\n")
			return string(line)
		default:
			fmt.Print(string(next_rune))
			line += string(next_rune)
		}

	}
}
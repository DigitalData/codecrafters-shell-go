package main

import (
	"log"
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

type CommandRune int
const (
	KeyCtrlC 		CommandRune = 3
	KeyTab 			CommandRune = 9
	KeyCtrlJ 		CommandRune = 10
	KeyEnter 		CommandRune = 13
	KeyBackspace 	CommandRune = 127
)

func parse_args(raw_line string) ([]string, *Outputs, error) {
	var args []string
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	set_output := UnsetOutput
	raw_line = strings.TrimSpace(raw_line)
	outputs := &Outputs{out_writer: os.Stdout, err_writer: os.Stdout}

	for _, r := range raw_line {
		quote := single_quotes || double_quotes
		if (!backslash) {
			continue_loop := false
			switch r {
			case '\\':
				if (set_output == UnsetOutput && !single_quotes) {
					backslash = true
					continue_loop = true
				}
			case '\'':
				if (set_output == UnsetOutput && !double_quotes) {
					single_quotes = !single_quotes
					continue_loop = true
				}
			case '"':
				if (set_output == UnsetOutput && !single_quotes) {
					double_quotes = !double_quotes
					continue_loop = true
				}
			case '>':
				if (!quote && !backslash) {
					switch set_output {
					case SetOutputOut:
						set_output = SetOutputOutAppend
					case SetOutputErr:
						set_output = SetOutputErrAppend
					default:
						if (len(current_arg) == 1 && current_arg[0] == '2') {
							set_output = SetOutputErr
						} else {
							set_output = SetOutputOut
						}
					}
					current_arg = ""
					continue_loop = true
				}
			default:
				if (!quote && unicode.IsSpace(r)) {
					continue_loop = true
					
					if (len(current_arg) == 0) {
						break
					} else if(set_output != UnsetOutput) {
						var err error
						err = outputs.update(current_arg, set_output)

						if (err != nil) {
							log.Fatal(err)
							return nil, nil, err
						}
						set_output = UnsetOutput
					} else {
						args = append(args, current_arg)
					}
					current_arg = ""
				}
			}

			if (continue_loop) {
				continue
			}
		}

		current_arg += string(r)
		backslash = false
	}

	if (set_output != UnsetOutput) {
		var err error
		err = outputs.update(current_arg, set_output)
		if (err != nil) {
			return nil, nil, err
		}
		set_output = UnsetOutput
	} else if (len(current_arg) > 0) {
		args = append(args, current_arg)
	}
	return args, outputs, nil
}

func loop(term_state *term.State) bool {
	var raw_line string = read_line()
	raw_line = strings.TrimSpace(raw_line)
	
	if len(raw_line) == 0 { return true }
	var args []string
	var outputs *Outputs
	var err error
	args, outputs, err = parse_args(raw_line)
	if (err != nil) {
		return true
	}
	if len(args) == 0 { return true }
	var cmd string = args[0]
	var cmd_args []string = args[1:]
	var has_args bool = len(cmd_args) > 0
	if cmd == CMD_EXIT { return false }

	var handler CMDHandler = handle_unknown
	switch cmd {
		case CMD_ECHO: 
			handler = handle_echo
		case CMD_TYPE:
			handler = handle_type
		case CMD_PWD:
			handler = handle_pwd
		case CMD_CD:
			handler = handle_cd
	}

	term.Restore(int(os.Stdin.Fd()), term_state)
	if handler != nil {
		handler(raw_line, cmd, cmd_args, has_args, outputs)
	}
	new_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	*term_state = *new_state

	return true
}

func main() {
	var term_state *term.State
	var err error
	term_state, err = term.MakeRaw(int(os.Stdin.Fd()))
	if (err != nil) {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), term_state)

	for loop(term_state) {}
}
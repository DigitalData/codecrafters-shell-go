package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"
)

const DEBUG = false
const CMD_EXIT = "exit"
// type CMDHandler func(raw_line string, cmd string, raw_args string, has_args bool)
type CMDHandler func(raw_line string, cmd string, cmd_args []string, has_args bool)

func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}
	
	var prog *exec.Cmd
	if has_args {
		prog = exec.Command(cmd, cmd_args...)
	} else {
		prog = exec.Command(cmd)
	}
	var std_out_err []byte
	std_out_err, err = prog.CombinedOutput()
	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Printf("%s", std_out_err)
}

const CMD_ECHO = "echo"
func handle_echo(_ string, _ string, cmd_args []string, _ bool) {
	var output string = strings.Join(cmd_args, " ")
	fmt.Printf("%s\n", output)
}

const CMD_TYPE = "type"
func handle_type(_ string, _ string, cmd_args []string, _ bool) {
	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD}
	for _, cmd_arg := range cmd_args {
		if slices.Contains(builtin_cmds, cmd_arg) {
			fmt.Printf("%s is a shell builtin\n", cmd_arg)
			continue
		} 
		
		var cmd_path string
		var err error
		cmd_path, err = exec.LookPath(cmd_arg)
		if err == nil {
			fmt.Printf("%s is %s\n", cmd_arg, cmd_path)
			continue
		} 
		
		fmt.Printf("%s: not found\n", cmd_arg)
	}
}

const CMD_PWD = "pwd"
func handle_pwd(_ string, _ string, _ []string, _ bool) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cwd)	
}

const CMD_CD = "cd"
func handle_cd(_ string, _ string, cmd_args []string, has_args bool) {
	var err error
	var home_dir string
	home_dir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var raw_args = "~"
	if has_args {
		raw_args = strings.Join(cmd_args, " ")
		raw_args = strings.ReplaceAll(raw_args, "~", home_dir)
	}
	
	err = os.Chdir(raw_args)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", raw_args)
	}
}

func parse_args(raw_line string) []string {
	var args []string
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	raw_line = strings.TrimSpace(raw_line)
	for _, r := range raw_line {

		quote := single_quotes || double_quotes

		if (!backslash) {
			continue_loop := false
			switch r {
			case '\\':
				if (!single_quotes) {
					backslash = true
					continue_loop = true
				}
			case '\'':
				if (!double_quotes) {
					single_quotes = !single_quotes
					continue_loop = true
				}
			case '"':
				if (!single_quotes) {
					double_quotes = !double_quotes
					continue_loop = true
				}
			default:
				if (unicode.IsSpace(r) && !quote) {
					if (len(current_arg) > 0) {
						args = append(args, current_arg)
						current_arg = ""
					}
					continue_loop = true
				}
			}

			if (continue_loop) {
				continue
			}
		}

		current_arg += string(r)
		backslash = false
	}

	if (len(current_arg) > 0) {
		args = append(args, current_arg)
	}
	return args
}


func loop() bool {
	fmt.Print("$ ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	raw_line := scanner.Text()
	raw_line = strings.TrimSpace(raw_line)
	
	if len(raw_line) == 0 { return true }
	var args []string  = parse_args(raw_line)
	if len(args) == 0 { return true }
	var cmd string = args[0]
	var cmd_args []string = args[1:]
	var has_args bool = len(cmd_args) > 0

	if DEBUG {
		fmt.Printf("DEBUG / cmd = \"%s\"\n", cmd)
		fmt.Printf("DEBUG / cmd_args = \"%v\"\n", cmd_args)
		fmt.Printf("DEBUG / has_args = %v\n", has_args)
	}

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

	if handler != nil {
		handler(raw_line, cmd, cmd_args, has_args)
	}

	return true
}


func main() {
	for loop() { }
}

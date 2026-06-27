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
type CMDHandler func(raw_line string, cmd string, raw_args string, has_args bool)

func handle_unknown(raw_line string, cmd string, raw_args string, has_args bool) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}
	
	var prog *exec.Cmd
	if has_args {
		var cmd_args []string = split_args(raw_args)
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
func handle_echo(_ string, _ string, raw_cmd_args string, _ bool) {
	var cmd_args []string = split_args(raw_cmd_args)
	var output string = strings.Join(cmd_args, " ")
	fmt.Printf("%s\n", output)
}

const CMD_TYPE = "type"
func handle_type(_ string, _ string, raw_args string, _ bool) {
	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD}
	cmd_args := split_args(raw_args)
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
func handle_pwd(_ string, _ string, _ string, _ bool) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cwd)	
}

const CMD_CD = "cd"
func handle_cd(_ string, _ string, raw_args string, has_args bool) {
	var err error
	var home_dir string
	home_dir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	
	if !has_args {
		raw_args = "~"
	}
	raw_args = strings.ReplaceAll(raw_args, "~", home_dir)
	
	err = os.Chdir(raw_args)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", raw_args)
	}
}

func split_args(raw_args string) []string {
	var args []string
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	raw_args = strings.TrimSpace(raw_args)
	for _, r := range raw_args {
		if (r == '\\' && !backslash) {
			backslash = true
			continue
		}

		if (r == '\'' && !backslash && !double_quotes) {
			single_quotes = !single_quotes
			continue
		}
		
		if (r == '"' && !backslash && !single_quotes) {
			double_quotes = !double_quotes
			continue
		}

		if (!backslash && !double_quotes && !single_quotes && unicode.IsSpace(r)) {
			if (len(current_arg) > 0) {
				args = append(args, current_arg)
				current_arg = ""
			}
			continue
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
	cmd, raw_args, has_args := strings.Cut(raw_line, " ")

	if DEBUG {
		fmt.Printf("DEBUG / cmd = \"%s\"\n", cmd)
		fmt.Printf("DEBUG / raw_cmd_args = \"%s\"\n", raw_args)
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
		handler(raw_line, cmd, raw_args, has_args)
	}

	return true
}


func main() {
	for loop() { }
}

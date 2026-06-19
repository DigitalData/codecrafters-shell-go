package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
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
	
	var cmd_args []string = strings.Split(raw_args, " ")
	var prog *exec.Cmd = exec.Command(cmd, cmd_args...)
	var std_out_err []byte
	std_out_err, err = prog.CombinedOutput()
	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Printf("%s", std_out_err)
}

const CMD_ECHO = "echo"
func handle_echo(_ string, _ string, raw_cmd_args string, _ bool) {
	fmt.Printf("%s\n", raw_cmd_args)
}

const CMD_TYPE = "type"
func handle_type(_ string, _ string, raw_args string, _ bool) {
	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD}
	cmd_args := strings.Split(raw_args, " ")
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
func handle_cd(_ string, _ string, raw_args string, _ bool) {
	var err error = os.Chdir(raw_args)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", raw_args)
	}
}


func loop() bool {
	fmt.Print("$ ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	raw_line := scanner.Text()

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

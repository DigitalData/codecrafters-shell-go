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
const CMD_ECHO = "echo"
const CMD_TYPE = "type"


func main() {

	for true {
		
		fmt.Print("$ ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		raw_args := scanner.Text()
	
		if len(raw_args) == 0 { continue }
		cmd, raw_cmd_args, has_args := strings.Cut(raw_args, " ")

		if DEBUG {
			fmt.Printf("DEBUG / cmd = \"%s\"\n", cmd)
			fmt.Printf("DEBUG / raw_cmd_args = \"%s\"\n", raw_cmd_args)
			fmt.Printf("DEBUG / has_args = %v\n", has_args)
		}

		if cmd == CMD_EXIT { break }

		switch cmd {
		case CMD_ECHO: 
			fmt.Printf("%s\n", raw_cmd_args)
		case CMD_TYPE:
			builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE}
			cmd_args := strings.Split(raw_cmd_args, " ")
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
		default:
			var err error
			_, err = exec.LookPath(cmd)
			if err != nil {
				fmt.Printf("%s: command not found\n", cmd)
				break
			}
			
			var cmd_args []string
			var prog *exec.Cmd
			var std_out_err []byte
			cmd_args = strings.Split(raw_cmd_args, " ")
			prog = exec.Command(cmd, cmd_args...)
			std_out_err, err = prog.CombinedOutput()
			if (err != nil) {
				log.Fatal(err)
			}
			fmt.Printf("%s", std_out_err)

		}

	}
	
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const CMD_EXIT = "exit"
type CMDHandler func(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *ShellIO)

func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *ShellIO) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		shell_io.errf("%s: command not found\n", cmd)
		return
	}

	var is_background bool
	is_background, cmd_args = extract_background_args(cmd_args)

	var prog *exec.Cmd
	if has_args {
		prog = exec.Command(cmd, cmd_args...)
	} else {
		prog = exec.Command(cmd)
	}
	
	prog.Stdin = shell_io.in_reader
	prog.Stdout = shell_io.out_writer
	prog.Stderr = shell_io.err_writer
	if (is_background) {
		var job_id, pid int
		job_id, pid, err = queue_job(raw_line, prog)
		fmt.Printf("[%d] %d\n", job_id, pid)
	} else {
		prog.Run()
	}
}

const CMD_ECHO = "echo"

func handle_echo(_ string, _ string, cmd_args []string, _ bool, shell_io *ShellIO) {
	var output string = strings.Join(cmd_args, " ")
	shell_io.outf("%s\n", output)
}

const CMD_TYPE = "type"

func handle_type(_ string, _ string, cmd_args []string, _ bool, shell_io *ShellIO) {
	builtin_cmds := []string{
		CMD_EXIT, 
		CMD_ECHO, 
		CMD_TYPE, 
		CMD_PWD,
		CMD_CD, 
		CMD_COMPLETE, 
		CMD_JOBS, 
		CMD_HISTORY,
	}
	for _, cmd_arg := range cmd_args {
		if slices.Contains(builtin_cmds, cmd_arg) {
			shell_io.outf("%s is a shell builtin\n", cmd_arg)
			continue
		}

		var cmd_path string
		var err error
		cmd_path, err = exec.LookPath(cmd_arg)
		if err == nil {
			shell_io.outf("%s is %s\n", cmd_arg, cmd_path)
			continue
		}

		shell_io.outf("%s: not found\n", cmd_arg)
	}
}

const CMD_PWD = "pwd"

func handle_pwd(_ string, _ string, _ []string, _ bool, shell_io *ShellIO) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	shell_io.outf("%s\n", cwd)
}

const CMD_CD = "cd"

func handle_cd(_ string, _ string, cmd_args []string, has_args bool, shell_io *ShellIO) {
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
		shell_io.outf("cd: %s: No such file or directory\n", raw_args)
	}
}

const CMD_COMPLETE = "complete"
var _completions map[string]string = make(map[string]string)

func handle_complete(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *ShellIO) {
	if (len(cmd_args) <= 1) {
		shell_io.err("complete: no flags given\n")
	}

	flag := cmd_args[0]
	num_args := len(cmd_args)

	switch flag {
	case "-C":
		if (num_args != 3) {
			shell_io.err("complete: expected two inputs for '-C' flag\n")
			return
		}
		program := cmd_args[num_args - 1]
		program_path := cmd_args[1]
		_completions[program] = program_path
	case "-r":
		if (num_args != 2) {
			shell_io.err("complete: expected two inputs for '-C' flag\n")
			return
		}
		program := cmd_args[num_args - 1]
		delete(_completions, program)
	case "-p":
		program := cmd_args[num_args - 1]
		program_path, exists := _completions[program]
		if (exists) {
			shell_io.outf("complete -C '%s' %s\n", program_path, program)
		} else {
			shell_io.errf("complete: %s: no completion specification\n", program)
		}
	default:
		shell_io.errf("complete: unsupported flag %s\n", flag)
	}
}

const CMD_JOBS = "jobs"

const CMD_HISTORY = "history"

func handle_history(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *ShellIO) {
	/* do nothing */
}
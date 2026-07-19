package main

import (
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const CMD_EXIT = "exit"
type CMDHandler func(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs)

func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		outputs.errf("%s: command not found\n", cmd)
		return
	}

	var prog *exec.Cmd
	if has_args {
		prog = exec.Command(cmd, cmd_args...)
	} else {
		prog = exec.Command(cmd)
	}

	prog.Stdout = outputs.out_writer
	prog.Stderr = outputs.err_writer
	prog.Run()
}

const CMD_ECHO = "echo"

func handle_echo(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
	var output string = strings.Join(cmd_args, " ")
	outputs.outf("%s\n", output)
}

const CMD_TYPE = "type"

func handle_type(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD, CMD_COMPLETE}
	for _, cmd_arg := range cmd_args {
		if slices.Contains(builtin_cmds, cmd_arg) {
			outputs.outf("%s is a shell builtin\n", cmd_arg)
			continue
		}

		var cmd_path string
		var err error
		cmd_path, err = exec.LookPath(cmd_arg)
		if err == nil {
			outputs.outf("%s is %s\n", cmd_arg, cmd_path)
			continue
		}

		outputs.outf("%s: not found\n", cmd_arg)
	}
}

const CMD_PWD = "pwd"

func handle_pwd(_ string, _ string, _ []string, _ bool, outputs *Outputs) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	outputs.outf("%s\n", cwd)
}

const CMD_CD = "cd"

func handle_cd(_ string, _ string, cmd_args []string, has_args bool, outputs *Outputs) {
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
		outputs.outf("cd: %s: No such file or directory\n", raw_args)
	}
}

const CMD_COMPLETE = "complete"


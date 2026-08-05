package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/shell_io"
)

func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool, pipe_io *shell_io.ShellIO) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		pipe_io.Errf("%s: command not found\n", cmd)
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
	
	prog.Stdin = pipe_io.InReader
	prog.Stdout = pipe_io.OutWriter
	prog.Stderr = pipe_io.ErrWriter
	if (is_background) {
		var job_id, pid int
		job_id, pid, err = queue_job(raw_line, prog)
		fmt.Printf("[%d] %d\n", job_id, pid)
	} else {
		prog.Run()
	}
}

func handle_echo(_ string, _ string, cmd_args []string, _ bool, pipe_io *shell_io.ShellIO) {
	var output string = strings.Join(cmd_args, " ")
	pipe_io.Outf("%s\n", output)
}

func handle_type(_ string, _ string, cmd_args []string, _ bool, pipe_io *shell_io.ShellIO) {
	for _, cmd_arg := range cmd_args {
		if slices.Contains(builtins, cmd_arg) {
			pipe_io.Outf("%s is a shell builtin\n", cmd_arg)
			continue
		}

		var cmd_path string
		var err error
		cmd_path, err = exec.LookPath(cmd_arg)
		if err == nil {
			pipe_io.Outf("%s is %s\n", cmd_arg, cmd_path)
			continue
		}

		pipe_io.Outf("%s: not found\n", cmd_arg)
	}
}

func handle_pwd(_ string, _ string, _ []string, _ bool, pipe_io *shell_io.ShellIO) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pipe_io.Outf("%s\n", cwd)
}

func handle_cd(_ string, _ string, cmd_args []string, has_args bool, pipe_io *shell_io.ShellIO) {
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
		pipe_io.Outf("cd: %s: No such file or directory\n", raw_args)
	}
}

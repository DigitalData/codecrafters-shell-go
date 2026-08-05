package command

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shell_io"
)

type Builtin string

const (
	BuiltinExit     Builtin = "exit"
	BuiltinEcho     Builtin = "echo"
	BuiltinType     Builtin = "type"
	BuiltinPWD      Builtin = "pwd"
	BuiltinCD       Builtin = "cd"
	BuiltinComplete Builtin = "complete"
	BuiltinJobs     Builtin = "jobs"
	BuiltinHistory  Builtin = "history"
)

type BuiltinHandler func(
	raw_line string, 
	cmd string, cmd_args []string, has_args bool, shell_io *shell_io.ShellIO,)

var builtin_map map[Builtin]BuiltinHandler = map[Builtin]BuiltinHandler{
	BuiltinExit: handle_exit,
	BuiltinEcho: handle_echo,
	BuiltinType: handle_type,
	BuiltinPWD: handle_pwd,
	BuiltinCD: handle_cd,
	BuiltinComplete: handle_complete,
	BuiltinJobs: handle_jobs,
	BuiltinHistory: handle_history,
}

var builtins []string = []string {
	string(BuiltinExit), 
	string(BuiltinEcho), 
	string(BuiltinType), 
	string(BuiltinPWD),
	string(BuiltinCD), 
	string(BuiltinComplete), 
	string(BuiltinJobs), 
	string(BuiltinHistory),
}

func GetHandler(cmd string) BuiltinHandler {
	handler, exists := builtin_map[Builtin(cmd)]
	if (!exists) {
		return handle_unknown
	}
	return handler
}

func handle_exit(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *shell_io.ShellIO) {
	os.Exit(0)
}
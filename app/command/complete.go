package command

import "github.com/codecrafters-io/shell-starter-go/app/shell_io"

var _completions map[string]string = make(map[string]string)

func GetCompleter(key string) (completer string, exists bool) {
	completer, exists = _completions[key]
	return completer, exists
}

func handle_complete(raw_line string, cmd string, cmd_args []string, has_args bool, pipe_io *shell_io.ShellIO) {
	if len(cmd_args) <= 1 {
		pipe_io.Err("complete: no flags given\n")
	}

	flag := cmd_args[0]
	num_args := len(cmd_args)

	switch flag {
	case "-C":
		if num_args != 3 {
			pipe_io.Err("complete: expected two inputs for '-C' flag\n")
			return
		}
		program := cmd_args[num_args-1]
		program_path := cmd_args[1]
		_completions[program] = program_path
	case "-r":
		if num_args != 2 {
			pipe_io.Err("complete: expected two inputs for '-C' flag\n")
			return
		}
		program := cmd_args[num_args-1]
		delete(_completions, program)
	case "-p":
		program := cmd_args[num_args-1]
		program_path, exists := _completions[program]
		if exists {
			pipe_io.Outf("complete -C '%s' %s\n", program_path, program)
		} else {
			pipe_io.Errf("complete: %s: no completion specification\n", program)
		}
	default:
		pipe_io.Errf("complete: unsupported flag %s\n", flag)
	}
}
package command

import "github.com/codecrafters-io/shell-starter-go/app/shell_io"

var _history []string

func AppendHistory(raw_line string) {
	_history = append(_history, raw_line)
}

func handle_history(raw_line string, cmd string, cmd_args []string, has_args bool, pipe_io *shell_io.ShellIO) {
	for idx, line := range _history {
		pipe_io.Outf("%5d  %s\n", idx+1, line)
	}
}
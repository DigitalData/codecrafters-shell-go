package main

var _history []string

func append_history(raw_line string) {
	_history = append(_history, raw_line)
}

func handle_history(raw_line string, cmd string, cmd_args []string, has_args bool, shell_io *ShellIO) {
	for idx, line := range _history {
		shell_io.outf("%5d  %s\n", idx+1, line)
	}
}
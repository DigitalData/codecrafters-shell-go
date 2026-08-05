package main

import (
	"os"
	"strings"

	"golang.org/x/term"
)


func run_pipeline(pipeline *ShellPipeline, next_pipeline *ShellPipeline) {

	if len(pipeline.args) == 0 { 
		return
	}
	
	if (next_pipeline != nil) {
		next_pipeline.shell_io.input(pipeline.shell_io.out_reader)
	}

	var cmd string = pipeline.args[0]
	var cmd_args []string = pipeline.args[1:]
	var has_args bool = len(cmd_args) > 0
	if cmd == CMD_EXIT { os.Exit(0) }

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
		case CMD_COMPLETE:
			handler = handle_complete
		case CMD_JOBS:
			handler = handle_jobs
		case CMD_HISTORY:
			handler = handle_history
	}

	if (handler != nil) {
		if (next_pipeline != nil) {
			go func () {
				defer pipeline.shell_io.out_writer.Close()
				defer pipeline.shell_io.err_writer.Close()
				handler(pipeline.raw, cmd, cmd_args, has_args, pipeline.shell_io)
			}()
		} else {
			handler(pipeline.raw, cmd, cmd_args, has_args, pipeline.shell_io)
		}
	}
}


func loop(term_state *term.State) bool {
	var raw_line string = read_line()
	raw_line = strings.TrimSpace(raw_line)
	
	if len(raw_line) == 0 { return true }
	append_history(raw_line)
	var pipelines []*ShellPipeline
	var err error
	pipelines, err = parse_args(raw_line)
	if (err != nil) {
		return true
	}

	term.Restore(int(os.Stdin.Fd()), term_state)

	num_pipelines := len(pipelines)
	if (num_pipelines == 0) { return true }
	for pipe_idx := range num_pipelines - 1 {
		pipeline := pipelines[pipe_idx]
		next_pipeline := pipelines[pipe_idx + 1]
		run_pipeline(pipeline, next_pipeline)
	}	
	run_pipeline(pipelines[num_pipelines - 1], nil)	
	print_and_reap_jobs(true, default_shell_io())

	new_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	*term_state = *new_state

	return true
}

func main() {
	var term_state *term.State
	var err error
	term_state, err = term.MakeRaw(int(os.Stdin.Fd()))
	if (err != nil) {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), term_state)
	
	for loop(term_state) {}
}
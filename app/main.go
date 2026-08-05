package main

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/command"
	"github.com/codecrafters-io/shell-starter-go/app/shell_io"
	"golang.org/x/term"
)


func run_pipeline(pipeline *ShellPipeline, next_pipeline *ShellPipeline) {

	if len(pipeline.args) == 0 { 
		return
	}
	
	if (next_pipeline != nil) {
		next_pipeline.pipe_io.Input(pipeline.pipe_io.OutReader)
	}

	var cmd string = pipeline.args[0]
	var cmd_args []string = pipeline.args[1:]
	var has_args bool = len(cmd_args) > 0

	var handler command.BuiltinHandler = command.GetHandler(cmd)
	if (handler != nil) {
		if (next_pipeline != nil) {
			go func () {
				defer pipeline.pipe_io.OutWriter.Close()
				defer pipeline.pipe_io.ErrWriter.Close()
				handler(pipeline.raw, cmd, cmd_args, has_args, pipeline.pipe_io)
			}()
		} else {
			handler(pipeline.raw, cmd, cmd_args, has_args, pipeline.pipe_io)
		}
	}
}


func loop(term_state *term.State) bool {
	var raw_line string = read_line()
	raw_line = strings.TrimSpace(raw_line)
	
	if len(raw_line) == 0 { return true }
	command.AppendHistory(raw_line)
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
	command.PrintAndReapJobs(true, shell_io.DefaultShellIO())

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
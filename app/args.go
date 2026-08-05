package main

import (
	"log"
	"strings"
	"unicode"

	"github.com/codecrafters-io/shell-starter-go/app/shell_io"
)

type ShellPipeline struct {
	raw      string
	args     []string
	pipe_io  *shell_io.ShellIO
}

func unparsed_pipeline() *ShellPipeline {
	return &ShellPipeline{"", []string{}, shell_io.PipelineShellIO(),}
}

func parse_args(raw_line string) (pipelines []*ShellPipeline, err error) {
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	set_output := shell_io.UnsetOutput
	raw_line = strings.TrimSpace(raw_line)
	pipeline := unparsed_pipeline()
	line := ""
	var pipe_io *shell_io.ShellIO = shell_io.DefaultShellIO()

	for _, r := range raw_line {
		line += string(r)
		quote := single_quotes || double_quotes
		if !backslash {
			continue_loop := false
			switch r {
			case '\\':
				if set_output == shell_io.UnsetOutput && !single_quotes {
					backslash = true
					continue_loop = true
				}
			case '\'':
				if set_output == shell_io.UnsetOutput && !double_quotes {
					single_quotes = !single_quotes
					continue_loop = true
				}
			case '"':
				if set_output == shell_io.UnsetOutput && !single_quotes {
					double_quotes = !double_quotes
					continue_loop = true
				}
			case '>':
				if !quote && !backslash {
					switch set_output {
					case shell_io.SetOutputOut:
						set_output = shell_io.SetOutputOutAppend
					case shell_io.SetOutputErr:
						set_output = shell_io.SetOutputErrAppend
					default:
						if len(current_arg) == 1 && current_arg[0] == '2' {
							set_output = shell_io.SetOutputErr
						} else {
							set_output = shell_io.SetOutputOut
						}
					}
					current_arg = ""
					continue_loop = true
				}
			case '|':
				if (!quote && !backslash) {
					pipeline.raw = strings.TrimRight(line, " |")
					pipelines = append(pipelines, pipeline)
					pipeline = unparsed_pipeline()
					line = ""
					pipe_io = shell_io.DefaultShellIO()
					continue_loop = true
				}
			default:
				if !quote && unicode.IsSpace(r) {
					continue_loop = true

					if len(current_arg) == 0 {
						break
					} else if set_output != shell_io.UnsetOutput {
						err = pipe_io.Update(current_arg, set_output)

						if err != nil {
							log.Fatal(err)
							return nil, err
						}
						set_output = shell_io.UnsetOutput
					} else {
						pipeline.args = append(pipeline.args, current_arg)
					}
					current_arg = ""
				}
			}

			if continue_loop {
				continue
			}
		}

		current_arg += string(r)
		backslash = false
	}

	if set_output != shell_io.UnsetOutput {
		err = pipe_io.Update(current_arg, set_output)
		if err != nil {
			return nil, err
		}
		set_output = shell_io.UnsetOutput
	} else if len(current_arg) > 0 {
		pipeline.args = append(pipeline.args, current_arg)
	}
	pipeline.pipe_io = pipe_io
	if (len(pipeline.args) > 0) {
		pipeline.raw = strings.TrimRight(line, " |")
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}
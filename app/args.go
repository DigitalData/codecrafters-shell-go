package main

import (
	"log"
	"strings"
	"unicode"
)

type ShellPipeline struct {
	raw      string
	args     []string
	shell_io  *ShellIO
}

func unparsed_pipeline() *ShellPipeline {
	return &ShellPipeline{"", []string{}, pipeline_shell_io(),}
}

func parse_args(raw_line string) (pipelines []*ShellPipeline, err error) {
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	set_output := UnsetOutput
	raw_line = strings.TrimSpace(raw_line)
	pipeline := unparsed_pipeline()
	line := ""
	var shell_io *ShellIO = default_shell_io()

	for _, r := range raw_line {
		line += string(r)
		quote := single_quotes || double_quotes
		if !backslash {
			continue_loop := false
			switch r {
			case '\\':
				if set_output == UnsetOutput && !single_quotes {
					backslash = true
					continue_loop = true
				}
			case '\'':
				if set_output == UnsetOutput && !double_quotes {
					single_quotes = !single_quotes
					continue_loop = true
				}
			case '"':
				if set_output == UnsetOutput && !single_quotes {
					double_quotes = !double_quotes
					continue_loop = true
				}
			case '>':
				if !quote && !backslash {
					switch set_output {
					case SetOutputOut:
						set_output = SetOutputOutAppend
					case SetOutputErr:
						set_output = SetOutputErrAppend
					default:
						if len(current_arg) == 1 && current_arg[0] == '2' {
							set_output = SetOutputErr
						} else {
							set_output = SetOutputOut
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
					shell_io = default_shell_io()
					continue_loop = true
				}
			default:
				if !quote && unicode.IsSpace(r) {
					continue_loop = true

					if len(current_arg) == 0 {
						break
					} else if set_output != UnsetOutput {
						err = shell_io.update(current_arg, set_output)

						if err != nil {
							log.Fatal(err)
							return nil, err
						}
						set_output = UnsetOutput
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

	if set_output != UnsetOutput {
		err = shell_io.update(current_arg, set_output)
		if err != nil {
			return nil, err
		}
		set_output = UnsetOutput
	} else if len(current_arg) > 0 {
		pipeline.args = append(pipeline.args, current_arg)
	}
	pipeline.shell_io = shell_io
	if (len(pipeline.args) > 0) {
		pipeline.raw = strings.TrimRight(line, " |")
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}
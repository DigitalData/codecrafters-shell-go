package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

// const CMD_EXIT = "exit"
// type CMDHandler func(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs)

type CommandRune int
const (
	KeyCtrlC 		CommandRune = 3
	KeyTab 			CommandRune = 9
	KeyCtrlJ 		CommandRune = 10
	KeyEnter 		CommandRune = 13
	KeyBackspace 	CommandRune = 127
)


type SetOutputMode int
const (
	UnsetOutput SetOutputMode = iota
	SetOutputOut
	SetOutputOutAppend
	SetOutputErr
	SetOutputErrAppend
)
type Outputs struct {
	out_writer io.Writer
	err_writer io.Writer
}
func (o *Outputs) outf(fstr string, vars ...any) {
	fmt.Fprintf(o.out_writer, fstr, vars...)
}
func (o *Outputs) out(str string) {
	fmt.Fprint(o.out_writer, str)
}
func (o *Outputs) errf(fstr string, vars ...any) {
	fmt.Fprintf(o.err_writer, fstr, vars...)
}
func (o *Outputs) err(str string) {
	fmt.Fprint(o.err_writer, str)
}

// type CRLFWriter struct {
// 	inner_writer io.Writer
// }
// func (cw *CRLFWriter) Write(raw_b []byte) (int, error) {
// 	var new_b []byte
// 	is_slash_r := false
// 	for _, b := range raw_b {
// 		if (!is_slash_r && string(b) == "\n") {
// 			new_b = append(new_b, []byte("\r\n")...)
// 			is_slash_r = false
// 			continue
// 		} else if (string(b) == "\r") {
// 			is_slash_r = true
// 		} else {
// 			is_slash_r = false
// 		}
// 		new_b = append(new_b, b)
// 	}
// 	return cw.inner_writer.Write(new_b)
// }

// func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs) {
// 	var err error
// 	_, err = exec.LookPath(cmd)
// 	if err != nil {
// 		outputs.errf("%s: command not found\r\n", cmd)
// 		return
// 	}

// 	var prog *exec.Cmd
// 	if has_args {
// 		prog = exec.Command(cmd, cmd_args...)
// 	} else {
// 		prog = exec.Command(cmd)
// 	}

// 	prog.Stdout = &CRLFWriter{outputs.out_writer}
// 	prog.Stderr = &CRLFWriter{outputs.err_writer}
// 	err = prog.Run()
// }

// const CMD_ECHO = "echo"
// func handle_echo(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
// 	var output string = strings.Join(cmd_args, " ")
// 	outputs.outf("%s\r\n", output)
// }

// const CMD_TYPE = "type"
// func handle_type(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
// 	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD}
// 	for _, cmd_arg := range cmd_args {
// 		if slices.Contains(builtin_cmds, cmd_arg) {
// 			outputs.outf("%s is a shell builtin\r\n", cmd_arg)
// 			continue
// 		} 
		
// 		var cmd_path string
// 		var err error
// 		cmd_path, err = exec.LookPath(cmd_arg)
// 		if err == nil {
// 			outputs.outf("%s is %s\r\n", cmd_arg, cmd_path)
// 			continue
// 		} 
		
// 		outputs.outf("%s: not found\r\n", cmd_arg)
// 	}
// }

// const CMD_PWD = "pwd"
// func handle_pwd(_ string, _ string, _ []string, _ bool, outputs *Outputs) {
// 	var cwd string
// 	var err error
// 	cwd, err = os.Getwd()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	outputs.outf("%s\r\n", cwd)	
// }

// const CMD_CD = "cd"
// func handle_cd(_ string, _ string, cmd_args []string, has_args bool, outputs *Outputs) {
// 	var err error
// 	var home_dir string
// 	home_dir, err = os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var raw_args = "~"
// 	if has_args {
// 		raw_args = strings.Join(cmd_args, " ")
// 		raw_args = strings.ReplaceAll(raw_args, "~", home_dir)
// 	}
	
// 	err = os.Chdir(raw_args)
// 	if err != nil {
// 		outputs.outf("cd: %s: No such file or directory\r\n", raw_args)
// 	}
// }

func get_filepath(filepath string) string {
	affixes := [...]string{"'", "\""}
	for _, r := range affixes {
		if (strings.HasPrefix(filepath, r) && strings.HasSuffix(filepath, r)) {
			filepath = strings.Trim(filepath, r)
		}
	}
	return filepath
}

func parse_args(raw_line string) ([]string, *Outputs, error) {
	var args []string
	current_arg := ""
	single_quotes := false
	double_quotes := false
	backslash := false
	set_output := UnsetOutput
	raw_line = strings.TrimSpace(raw_line)
	outputs := &Outputs{out_writer: os.Stdout, err_writer: os.Stdout}

	for _, r := range raw_line {
		quote := single_quotes || double_quotes
		if (!backslash) {
			continue_loop := false
			switch r {
			case '\\':
				if (set_output == UnsetOutput && !single_quotes) {
					backslash = true
					continue_loop = true
				}
			case '\'':
				if (set_output == UnsetOutput && !double_quotes) {
					single_quotes = !single_quotes
					continue_loop = true
				}
			case '"':
				if (set_output == UnsetOutput && !single_quotes) {
					double_quotes = !double_quotes
					continue_loop = true
				}
			case '>':
				if (!quote && !backslash) {
					switch set_output {
					case SetOutputOut:
						set_output = SetOutputOutAppend
					case SetOutputErr:
						set_output = SetOutputErrAppend
					default:
						if (len(current_arg) == 1 && current_arg[0] == '2') {
							set_output = SetOutputErr
						} else {
							set_output = SetOutputOut
						}
					}
					current_arg = ""
					continue_loop = true
				}
			default:
				if (!quote && unicode.IsSpace(r)) {
					continue_loop = true
					
					if (len(current_arg) == 0) {
						break
					} else if(set_output != UnsetOutput) {
						var err error
						var filepath string = get_filepath(current_arg)
						switch set_output {
						case SetOutputOut:
							outputs.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
						case SetOutputErr:
							outputs.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
						case SetOutputOutAppend:
							outputs.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_APPEND, 0644)
						case SetOutputErrAppend:
							outputs.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_APPEND, 0644)
						}
						if (err != nil) {
							return nil, nil, err
						}
						set_output = UnsetOutput
					} else {
						args = append(args, current_arg)
					}
					current_arg = ""
				}
			}

			if (continue_loop) {
				continue
			}
		}

		current_arg += string(r)
		backslash = false
	}

	if (set_output != UnsetOutput) {
		filepath := get_filepath(current_arg)
		var err error
		switch set_output {
		case SetOutputOut:
			outputs.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		case SetOutputErr:
			outputs.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		case SetOutputOutAppend:
			outputs.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		case SetOutputErrAppend:
			outputs.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		}
		if (err != nil) {
			return nil, nil, err
		}
		set_output = UnsetOutput
	} else if (len(current_arg) > 0) {
		args = append(args, current_arg)
	}
	return args, outputs, nil
}

func loop() bool {
	var raw_line string = read_line()
	raw_line = strings.TrimSpace(raw_line)
	
	if len(raw_line) == 0 { return true }
	var args []string
	var outputs *Outputs
	var err error
	args, outputs, err = parse_args(raw_line)
	if (err != nil) {
		return true
	}
	if len(args) == 0 { return true }
	var cmd string = args[0]
	var cmd_args []string = args[1:]
	var has_args bool = len(cmd_args) > 0
	if cmd == CMD_EXIT { return false }

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
	}

	if handler != nil {
		handler(raw_line, cmd, cmd_args, has_args, outputs)
	}

	return true
}

func main() {
	var old_state *term.State
	var err error
	old_state, err = term.MakeRaw(int(os.Stdin.Fd()))
	if (err != nil) {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), old_state)

	for loop() {}
}
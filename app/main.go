package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"
)

const CMD_EXIT = "exit"
// type CMDHandler func(raw_line string, cmd string, raw_args string, has_args bool)
type CMDHandler func(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs)

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


func handle_unknown(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs) {
	var err error
	_, err = exec.LookPath(cmd)
	if err != nil {
		// fmt.Printf("%s: command not found\n", cmd)
		outputs.errf("%s: command not found\n", cmd)
		return
	}
	
	var prog *exec.Cmd
	if has_args {
		prog = exec.Command(cmd, cmd_args...)
	} else {
		prog = exec.Command(cmd)
	}
	// var std_out_err []byte
	prog.Stdout = outputs.out_writer
	prog.Stderr = outputs.err_writer
	err = prog.Run()
}

const CMD_ECHO = "echo"
func handle_echo(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
	var output string = strings.Join(cmd_args, " ")
	// fmt.Printf("%s\n", output)
	outputs.outf("%s\n", output)
}

const CMD_TYPE = "type"
func handle_type(_ string, _ string, cmd_args []string, _ bool, outputs *Outputs) {
	builtin_cmds := []string{CMD_EXIT, CMD_ECHO, CMD_TYPE, CMD_PWD, CMD_CD}
	for _, cmd_arg := range cmd_args {
		if slices.Contains(builtin_cmds, cmd_arg) {
			// fmt.Printf("%s is a shell builtin\n", cmd_arg)
			outputs.outf("%s is a shell builtin\n", cmd_arg)
			continue
		} 
		
		var cmd_path string
		var err error
		cmd_path, err = exec.LookPath(cmd_arg)
		if err == nil {
			// fmt.Printf("%s is %s\n", cmd_arg, cmd_path)
			outputs.outf("%s is %s\n", cmd_arg, cmd_path)
			continue
		} 
		
		// fmt.Printf("%s: not found\n", cmd_arg)
		outputs.outf("%s: not found\n", cmd_arg)
	}
}

const CMD_PWD = "pwd"
func handle_pwd(_ string, _ string, _ []string, _ bool, outputs *Outputs) {
	var cwd string
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", cwd)	
	outputs.outf("%s\n", cwd)	
}

const CMD_CD = "cd"
func handle_cd(_ string, _ string, cmd_args []string, has_args bool, outputs *Outputs) {
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
		// fmt.Printf("cd: %s: No such file or directory\n", raw_args)
		outputs.outf("cd: %s: No such file or directory\n", raw_args)
	}
}

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
						// case SetOutputErr:
						// 	outputs.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
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
		// case SetOutputErr:
		// 	outputs.err_writer, err = os.OpenFile(currfilepathent_arg, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		}
		if (err != nil) {
			fmt.Println(err)
			return nil, nil, err
		}
		set_output = UnsetOutput
	} else if (len(current_arg) > 0) {
		args = append(args, current_arg)
	}
	return args, outputs, nil
}

func loop() bool {
	fmt.Print("$ ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	raw_line := scanner.Text()
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
	for loop() { }
}

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type SetOutputMode int

const (
	UnsetOutput SetOutputMode = iota
	SetOutputOut
	SetOutputOutAppend
	SetOutputErr
	SetOutputErrAppend
)

type ShellIO struct {
	in_reader io.ReadCloser
	out_writer io.WriteCloser
	out_reader io.ReadCloser
	err_writer io.WriteCloser
	err_reader io.ReadCloser
}

func (o *ShellIO) outf(fstr string, vars ...any) {
	fmt.Fprintf(o.out_writer, fstr, vars...)
}
func (o *ShellIO) out(str string) {
	fmt.Fprint(o.out_writer, str)
}
func (o *ShellIO) errf(fstr string, vars ...any) {
	fmt.Fprintf(o.err_writer, fstr, vars...)
}
func (o *ShellIO) err(str string) {
	fmt.Fprint(o.err_writer, str)
}
func (o *ShellIO) input(in_reader io.ReadCloser) {
	o.in_reader = in_reader
}
func (o *ShellIO) update(arg string, set_output SetOutputMode) (err error) {
	var read_writer io.ReadWriteCloser
	var filepath string = get_filepath(arg)
	switch set_output {
	case SetOutputOut:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		o.out_writer = read_writer
		o.out_reader = read_writer
	case SetOutputErr:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		o.err_writer = read_writer
		o.err_reader = read_writer
	case SetOutputOutAppend:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		o.out_writer = read_writer
		o.out_reader = read_writer
	case SetOutputErrAppend:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		o.err_writer = read_writer
		o.err_reader = read_writer
	}
	return err
}

func default_shell_io() *ShellIO {
	return &ShellIO{
		in_reader: os.Stdin,
		out_writer: os.Stdout, 
		out_reader: os.Stdout,
		err_writer: os.Stdout,
		err_reader: os.Stdout,
	}
}

func pipeline_shell_io() *ShellIO {
	reader, writer := io.Pipe()
	return &ShellIO{
		in_reader: os.Stdin,
		out_writer: writer, 
		out_reader: reader,
		err_writer: writer,
		err_reader: reader,
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
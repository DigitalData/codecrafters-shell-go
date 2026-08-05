package shell_io

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
	InReader io.ReadCloser
	OutWriter io.WriteCloser
	OutReader io.ReadCloser
	ErrWriter io.WriteCloser
	ErrReader io.ReadCloser
}

func (o *ShellIO) Outf(fstr string, vars ...any) {
	fmt.Fprintf(o.OutWriter, fstr, vars...)
}
func (o *ShellIO) Out(str string) {
	fmt.Fprint(o.OutWriter, str)
}
func (o *ShellIO) Errf(fstr string, vars ...any) {
	fmt.Fprintf(o.ErrWriter, fstr, vars...)
}
func (o *ShellIO) Err(str string) {
	fmt.Fprint(o.ErrWriter, str)
}
func (o *ShellIO) Input(in_reader io.ReadCloser) {
	o.InReader = in_reader
}
func (o *ShellIO) Update(arg string, set_output SetOutputMode) (err error) {
	var read_writer io.ReadWriteCloser
	var filepath string = get_filepath(arg)
	switch set_output {
	case SetOutputOut:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		o.OutWriter = read_writer
		o.OutReader = read_writer
	case SetOutputErr:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
		o.ErrWriter = read_writer
		o.ErrReader = read_writer
	case SetOutputOutAppend:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		o.OutWriter = read_writer
		o.OutReader = read_writer
	case SetOutputErrAppend:
		read_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		o.ErrWriter = read_writer
		o.ErrReader = read_writer
	}
	return err
}

func DefaultShellIO() *ShellIO {
	return &ShellIO{
		InReader: os.Stdin,
		OutWriter: os.Stdout, 
		OutReader: os.Stdout,
		ErrWriter: os.Stdout,
		ErrReader: os.Stdout,
	}
}

func PipelineShellIO() *ShellIO {
	reader, writer := io.Pipe()
	return &ShellIO{
		InReader: os.Stdin,
		OutWriter: writer, 
		OutReader: reader,
		ErrWriter: writer,
		ErrReader: reader,
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
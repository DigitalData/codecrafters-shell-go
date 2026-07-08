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
func (o *Outputs) update(arg string, set_output SetOutputMode) (err error) {
	var filepath string = get_filepath(arg)
	switch set_output {
	case SetOutputOut:
		o.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	case SetOutputErr:
		o.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	case SetOutputOutAppend:
		o.out_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	case SetOutputErrAppend:
		o.err_writer, err = os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	}
	return err
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
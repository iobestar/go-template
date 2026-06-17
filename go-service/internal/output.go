package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type output struct {
	writer io.Writer
}

func (o *output) Write(p []byte) (n int, err error) {
	return o.writer.Write(p)
}

func newOutput(w io.Writer) *output {
	return &output{
		writer: w,
	}
}

func newStdout() *output {
	return newOutput(os.Stdout)
}

func newStderr() *output {
	return newOutput(os.Stderr)
}

type writeFunc func(w io.Writer, a ...any) (n int, err error)

func (o *output) write(value any, w writeFunc) error {
	if nil == value {
		return nil
	}

	if err, ok := value.(error); ok {
		_, err := w(o.writer, fmt.Sprintf("ERROR: %s", err.Error()))
		return err
	}

	_, err := w(o.writer, fmt.Sprintf("%v", value))
	return err
}
func Write(value interface{}) error {
	return newStdout().write(value, fmt.Fprint)
}

func MustWrite(value interface{}) {
	if err := Write(value); nil != err {
		Fatal(err, 1)
	}
}

func Writeln(value interface{}) error {
	return newStdout().write(value, fmt.Fprintln)
}

func MustWriteln(value interface{}) {
	if err := Writeln(value); nil != err {
		Fatal(err, 1)
	}
}

func (o *output) writeFormat(w writeFunc, value any) error {
	if nil == value {
		return nil
	}

	b, err := json.Marshal(value)
	if nil != err {
		return err
	}

	var indent bytes.Buffer
	if err := json.Indent(&indent, b, "", "  "); nil != err {
		return err
	}

	_, err = w(o.writer, string(indent.Bytes()))
	return err
}

func WriteFormat(value any) {
	if err := newStdout().writeFormat(fmt.Fprintln, value); nil != err {
		Fatal(err, 1)
	}
}

func Fatal(err error, code int) {
	if err := newStdout().write(err, fmt.Fprintln); nil != err {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("FATAL: %s", err.Error()))
	}
	os.Exit(code)
}

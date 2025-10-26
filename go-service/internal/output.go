package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

type output struct {
	writer io.Writer
	format string
}

func newStdOutput() *output {
	return newOutput(os.Stdout)
}

func newOutput(writer io.Writer) *output {
	return &output{
		writer: writer,
	}
}

type writeFunc func(w io.Writer, a ...any) (n int, err error)

func (o *output) write(w writeFunc, value any) error {
	if nil == value {
		return nil
	}

	if err, ok := value.(error); ok {
		var exitError *ExitError
		if errors.As(err, &exitError) {
			if _, err := w(o.writer, fmt.Sprintf("ERROR: %s", err.Error())); nil != err {
				return err
			}
			os.Exit(exitError.exitCode)
		}

		if _, err := w(o.writer, fmt.Sprintf("ERROR: %s", err.Error())); nil != err {
			return err
		}
		return nil
	}

	_, err := w(o.writer, fmt.Sprintf("%v", value))
	return err
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

func Write(value any) {
	if err := newStdOutput().write(fmt.Fprint, value); nil != err {
		FatalError(err)
	}
}

func Writeln(value any) {
	if err := newStdOutput().write(fmt.Fprintln, value); nil != err {
		FatalError(NewExitError(err, 1))
	}
}

func WriteFormat(value any) {
	if err := newStdOutput().writeFormat(fmt.Fprintln, value); nil != err {
		FatalError(NewExitError(err, 1))
	}
}

func FatalError(err error) {
	if _, err := fmt.Fprintf(os.Stdout, "ERROR: %s", err.Error()); nil != err {
		os.Exit(1)
	}

	var exitError *ExitError
	if errors.As(err, &exitError) {
		os.Exit(exitError.exitCode)
	}
	os.Exit(1)
}

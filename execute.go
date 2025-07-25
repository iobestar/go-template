package template

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
	Env  []string
	Dir  string
}

func Execute(ctx context.Context, cmd Command) (int, []byte, []byte, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code, err := ExecuteWithWriter(ctx, cmd, stdout, stderr)
	return code, stdout.Bytes(), stderr.Bytes(), err
}

func ExecuteWithWriter(ctx context.Context, cmd Command, stdout io.Writer, stderr io.Writer) (int, error) {
	command := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	command.Dir = cmd.Dir
	command.Env = cmd.Env
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Start(); nil != err {
		return 0, errors.Wrap(err, "unable to start command")
	}

	if err := command.Wait(); nil != err {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, errors.New("execution timeout")
		}
		var e *exec.ExitError
		switch {
		case errors.As(err, &e):
			return e.ExitCode(), nil
		default:
			return 0, errors.Wrap(err, "unexpected execution error")
		}
	}
	return 0, nil
}

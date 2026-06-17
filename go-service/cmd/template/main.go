package main

import (
	"errors"

	"github.com/iobestar/go-template/go-service/cmd"
	"github.com/iobestar/go-template/go-service/internal"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var exitError *internal.ExitCodeError
		if ok := errors.As(err, &exitError); ok {
			internal.Fatal(exitError, exitError.ExitCode())
		} else {
			internal.Fatal(err, 1)
		}
	}
}

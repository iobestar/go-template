package main

import (
	"github.com/iobestar/go-template/go-service/cmd"
	"github.com/iobestar/go-template/go-service/internal"
)

func main() {
	if err := cmd.Execute(); nil != err {
		internal.FatalError(err)
	}
}

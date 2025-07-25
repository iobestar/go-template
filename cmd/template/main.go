package main

import (
	"github.com/iobestar/go-template/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); nil != err {
		log.Fatalln(err.Error())
	}
}

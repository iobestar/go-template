package template

import (
	"fmt"
	"strings"

	"github.com/mattn/go-tty"
	"github.com/pkg/errors"
)

type InputOptions struct {
	Default  string
	NoEcho   bool
	Required bool
}

func ReadInput(prompt string, options InputOptions) (string, error) {
	t, err := tty.Open()
	if nil != err {
		return "", err
	}

	var read func() (string, error)
	if options.NoEcho {
		read = t.ReadPasswordNoEcho
	} else {
		read = t.ReadString
	}

	if len(options.Default) > 0 {
		prompt = fmt.Sprintf("%s [%s]: ", prompt, options.Default)
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	for i := 3; i > 0; i-- {
		fmt.Print(prompt)
		var input string

		input, err := read()
		if nil != err {
			return "", err
		}
		input = strings.TrimSpace(input)
		if len(input) > 0 {
			return input, nil
		}

		if len(options.Default) > 0 {
			return options.Default, nil
		}

		if !options.Required {
			return "", nil
		}

		attempt := i - 1
		if attempt > 0 {
			fmt.Println(fmt.Sprintf("input required - try again (%d)", attempt))
		}
	}
	return "", errors.New("empty input not valid")
}

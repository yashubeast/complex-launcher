package internal

import (
	"cl/sources"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

type ExecuteFunc func(mappedValue string, input string) error

var executeMap = map[string]ExecuteFunc{
	"url": executeUrl,
	"copy": executeCopy,
	"copy_and_notify": executeCopyAndNotify,
}

func executeUrl(mappedValue string, input string) error {
	target := fmt.Sprintf(
		mappedValue,
		url.QueryEscape(input),
	)

	return exec.Command(
		"xdg-open",
		target,
	).Start()
}

func executePrefix(item sources.Item) error {
	execute, ok := executeMap[item.PrefixExecuteType]
	if !ok {
		return fmt.Errorf(
			"unknown execute function: %s",
			item.PrefixExecuteType,
		)
	}
	return execute(item.Cmd, item.PrefixInput)
}

func executeCopy(mappedValue string, input string) error {
	target := fmt.Sprintf(
		mappedValue,
		input,
	)

	// TODO: make this command alterable in config
	cmd := exec.Command(
		"xclip",
		"-selection",
		"clipboard",
	)

	cmd.Stdin = strings.NewReader(target)
	return cmd.Run()
}

func executeCopyAndNotify(mappedValue string, input string) error {
	target := fmt.Sprintf(
		mappedValue,
		input,
	)

	executeCopy(mappedValue, input)

	// TODO: also make this alterable
	return exec.Command(
		"notify-send",
		"cl copied",
		target,
	).Run()
}

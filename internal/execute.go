package internal

import (
	"cl/sources"
	"fmt"
	"net/url"
	"os/exec"
)

type ExecuteFunc func(mappedValue string, input string) error

var executeMap = map[string]ExecuteFunc{
	"url": executeUrl,
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

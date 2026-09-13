package sources

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
)

type Prefix struct {
	PrefixString string              `json:"prefix"`
	Kind         string              `json:"kind"`
	Items        []item              `json:"items"`
	GetItems     func(string) []Item `json:"-"`
	ExecuteType  string              `json:"execute_type"`
}

type item struct {
	Name string `json:"name"`
	Cmd string `json:"cmd"`
}

type PrefixSource struct {
	Prefixes []Prefix
	Query string
}

func (s *PrefixSource) List() ([]Item, error) {
	// loop through the given prefixes in the PrefixSource and filter valid ones according to input
	for _, prefix := range s.Prefixes {
		input, ok := strings.CutPrefix(
			s.Query,
			prefix.PrefixString,
		)
		if !ok { continue }

		// Prefix.Kind takes priority
		if prefix.Kind == "calculator" {
			return getCalculatorItems(input, prefix.ExecuteType)
		}
		// GetItems takes priority
		if prefix.GetItems != nil {
			return prefix.GetItems(input), nil
		}
		// otherwise generate items from ItemMap
		items := make([]Item, 0, len(prefix.Items))
		for _, item := range prefix.Items {
			items = append(items, Item{
				Name: item.Name,
				Cmd: item.Cmd,
				PrefixExecuteType: prefix.ExecuteType,
				PrefixInput: input,
			})
		}
		return items, nil
	}
	return nil, nil
}

func getCalculatorItems(input, executeType string) ([]Item, error) {
	dummyItemList := []Item{
		{ Name: "..." },
		{ Name: input + " = ..." },
	}
	input = strings.TrimSpace(input)
	if input == "" { return dummyItemList, nil }

	result, err := calculate(input)
	if err != nil { return dummyItemList, nil }

	return []Item{
		{
			Name:              result,
			Cmd:               "%s",
			PrefixExecuteType: executeType,
			PrefixInput:       result,
		},
		{
			Name:              input + " = " + result,
			Cmd:               "%s",
			PrefixExecuteType: executeType,
			PrefixInput:       input + " = " + result,
		},
	}, nil
}

func calculate(input string) (string, error) {
	program, err := expr.Compile(input)
	if err != nil { return "", err }

	value, err := expr.Run(program, nil)
	if err != nil { return "", err }

	return fmt.Sprint(value), nil
}

package sources

import "strings"

type Prefix struct {
	PrefixString string              `json:"prefix"`

	ItemMap      map[string]string   `json:"items"`
	GetItems     func(string) []Item `json:"-"`
	ExecuteType  string              `json:"execute_type"`
}
type PrefixConfig struct {
	Prefixes []Prefix `json:"prefixes"`
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

		// GetItems takes priority
		if prefix.GetItems != nil {
			return prefix.GetItems(input), nil
		}
		// otherwise generate items from ItemMap
		items := make([]Item, 0, len(prefix.ItemMap))
		for name, cmd := range prefix.ItemMap {
			items = append(items, Item{
				Name: name,
				Cmd: cmd,
				PrefixExecuteType: prefix.ExecuteType,
				PrefixInput: input,
			})
		}
		return items, nil
	}
	return nil, nil
}

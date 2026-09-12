package sources

type Item struct {
	Name              string
	Cmd               string
	PrefixExecuteType string
	PrefixInput       string
}

type Source interface {
	List() ([]Item, error)
}

func GetNames(items []Item) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}

func GetItems(names []string) []Item {
	items := make([]Item, len(names))
	for i, name := range names {
		items[i] = Item{
			Name: name,
			Cmd: name,
		}
	}
	return items
}

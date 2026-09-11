package sources

type Item struct {
	Name string
	Cmd  string
	Args []string
}

type Source interface {
	List() ([]Item, error)
}

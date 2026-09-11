package sources

import (
	"os"
	"path/filepath"
	"strings"
)

type PathSource struct{}

func (PathSource) List() ([]Item, error) {
	var items []Item
  seen := make(map[string]struct{})

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		entries, err := os.ReadDir(dir)
		if err != nil { continue }

		for _, entry := range entries {
			if entry.IsDir() { continue }

			name := entry.Name()
			if strings.Contains(name, " ") { continue }

			info, err := entry.Info()
			if err != nil || info.Mode()&0111 == 0 { continue }

			// first executable with this name wins, like PATH lookup
			if _, ok := seen[name]; ok { continue }
			seen[name] = struct{}{}

			items = append(items, Item{
				Name: name,
				Cmd: name,
			})
		}
	}
	return items, nil
}

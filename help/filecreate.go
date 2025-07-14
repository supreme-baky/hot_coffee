package help

import (
	"os"
	"path/filepath"
)

func CreateDataDirWithFiles(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	files := map[string]string{
		"inventory.json":  "[]",
		"menu_items.json": "[]",
		"orders.json":     "[]",
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}

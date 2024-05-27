package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/helper"
)

type Cache struct {
	Packages map[string]string
}

func NewCache() *Cache {
	c := &Cache{Packages: make(map[string]string)}
	if helper.HasHomeDir() {
		files, err := os.ReadDir(helper.HomeDir())
		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					c.Packages[file.Name()] = filepath.Join(helper.HomeDir(), file.Name())
				}
			}

			return c
		}
	}

	return nil
}

func (c *Cache) IsInCache(name, version string) bool {
	if _, yes := c.Packages[fmt.Sprintf("%s@%s", name, version)]; yes {
		return true
	}

	return false
}

func (c *Cache) IsEmpty() bool {
	return len(c.Packages) == 0
}

func (c *Cache) Path(name, version string) string {
	return filepath.Join(helper.HomeDir(), fmt.Sprintf("%s@%s", name, version))
}

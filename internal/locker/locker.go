package locker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/helper"
)

type Element struct {
	Version   string `json:"version"`
	Resolved  string `json:"Resolved"`
	Integrity string `json:"integrity"`
	Path      string `json:"path"`
}

type File struct {
	Elements map[string]*Element
}

func NewLocker() *File {
	lr := &File{
		Elements: make(map[string]*Element),
	}

	lr.Read()

	return lr
}

func (l *File) Add(i *Element, key string) {
	l.Elements[key] = i
}

func (l *File) Write() error {
	path := filepath.Join(helper.GetCurrentDirPath(), "ivar.lock")

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			fmt.Printf("error while deleting existing ivar.lock file: %s \n", err.Error())
		}
	}

	data, err := json.MarshalIndent(l.Elements, "", "	")
	if err != nil {
		return fmt.Errorf("error in encoding ivar.lock: %s", err.Error())
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("error in writing ivar.lock: %s", err.Error())
	}

	return nil
}

func (l *File) Read() error {
	path := filepath.Join(helper.GetCurrentDirPath(), "ivar.lock")

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error in reading ivar.lock: %s", err.Error())
		}

		var elements map[string]*Element

		err = json.Unmarshal(data, &elements)

		if err != nil {
			return fmt.Errorf("error in reading ivar.lock: %s", err.Error())
		}

		l.Elements = elements
	}

	return nil
}

func (l *File) GetVersion(name, version string) string {
	if l.Elements != nil && len(l.Elements) > 0 {
		key := fmt.Sprintf("%s@%s", name, version)
		if el, exists := l.Elements[key]; exists {
			return el.Version
		}
	}

	return ""
}

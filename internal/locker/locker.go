package locker

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/harshalslimaye/ivar/internal/graph"
)

type LockItem struct {
	Version   string `json:"version"`
	Resolved  string `json:"Resolved"`
	Integrity string `json:"integrity"`
	Path      string `json:"path"`
}

type Locker struct {
	Items map[string]*LockItem
}

func NewLocker() *Locker {
	return &Locker{
		Items: make(map[string]*LockItem),
	}
}

func (l *Locker) Add(n *graph.Node, path string) {
	key := fmt.Sprintf("%s@%s", n.Name(), n.Package.RawVersion)

	l.Items[key] = &LockItem{
		Version:   n.Version(),
		Resolved:  n.TarballUrl,
		Integrity: n.Integrity,
		Path:      path,
	}
}

func (l *Locker) Write() error {
	data, err := json.MarshalIndent(l.Items, "", "	")
	if err != nil {
		return fmt.Errorf("error in encoding ivar.lock: %s", err.Error())
	}
	err = os.WriteFile("ivar.lock", data, 0644)
	if err != nil {
		return fmt.Errorf("error in writing ivar.lock: %s", err.Error())
	}

	return nil
}

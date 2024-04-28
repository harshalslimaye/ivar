package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/logrusorgru/aurora"
)

type Dependency struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

var cache map[string]map[string]string
var cacheMutex sync.RWMutex

func init() {
	cache = make(map[string]map[string]string)
}

func FetchDependencies(name, version string) (map[string]string, error) {
	cacheMutex.RLock()
	if deps, ok := cache[name+version]; ok {
		cacheMutex.RUnlock()
		fmt.Printf(aurora.BrightYellow("%s@%s is chilling in the cache!\n").String(), name, version)
		return deps, nil
	}
	cacheMutex.RUnlock()

	endpoint := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	fmt.Println(aurora.BrightCyan("Fetching dependencies for " + name + "@" + version))
	var dep Dependency

	res, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch dependency %s@%s: %s", aurora.Yellow(name), aurora.Yellow(version), aurora.Yellow(res.Status))
	}

	err = json.NewDecoder(res.Body).Decode(&dep)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	cache[name+version] = dep.Dependencies
	cacheMutex.Unlock()

	return dep.Dependencies, nil
}

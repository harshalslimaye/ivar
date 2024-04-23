package graph

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DependencyList struct {
	Dependencies []*Dependency
	Visited      map[string]string
}

func BuildList(l *DependencyList, deps map[string]string) {
	for name, version := range deps {
		if _, ok := l.Visited[name]; !ok {
			FetchDependency(name, version, l)
		}
	}
}

func FetchDependency(name string, version string, list *DependencyList) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch dependency %s@%s: %v", name, version, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch dependency %s@%s: %s", name, version, res.Status)
	}

	var dep struct {
		Name         string
		Version      string
		Dependencies map[string]string
	}
	err = json.NewDecoder(res.Body).Decode(&dep)
	if err != nil {
		return fmt.Errorf("failed to decode response body for dependency %s@%s: %v", name, version, err)
	}
	list.Dependencies = append(list.Dependencies, &Dependency{Name: name, Version: version})
	list.Visited[name] = version

	if len(dep.Dependencies) > 0 {
		BuildList(list, dep.Dependencies)
	}

	return nil
}

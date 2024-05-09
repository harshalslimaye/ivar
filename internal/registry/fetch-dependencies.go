package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/logrusorgru/aurora"
)

type Dependency struct {
	Name         string
	Version      string
	Dependencies map[string]string
	Bin          map[string]string
}

func FetchDependencies(name, version string) (*Dependency, error) {
	endpoint := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	var dep Dependency

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch dependency %s@%s: %s", aurora.Yellow(name), aurora.Yellow(version), aurora.Yellow(res.Status))
	}

	if err := json.NewDecoder(res.Body).Decode(&dep); err != nil {
		return nil, err
	}

	return &dep, nil
}

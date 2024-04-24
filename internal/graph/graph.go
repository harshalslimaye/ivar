package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/Masterminds/semver/v3"
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
			GetAvailableVersions(name)
		}
	}
}

func GetVersion(name string, version string) string {
	return FindExactVersion(version, GetAvailableVersions(name))
}

// FindExactVersion finds the exact version that satisfies the given constraint
func FindExactVersion(constraint string, versions []*semver.Version) string {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		// Handle invalid constraint
		fmt.Printf("Error parsing version constraint %s: %v\n", constraint, err)
		return ""
	}

	var compatibleVersions []*semver.Version
	for _, v := range versions {
		if c.Check(v) {
			compatibleVersions = append(compatibleVersions, v)
		}
	}

	if len(compatibleVersions) == 0 {
		// No compatible versions found
		return ""
	}

	// Sort compatible versions
	sort.Sort(semver.Collection(compatibleVersions))

	// Return the highest compatible version
	return compatibleVersions[len(compatibleVersions)-1].String()
}

func GetAvailableVersions(name string) []*semver.Version {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Parse the JSON response
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Extract keys of the "versions" array
	versions := data["versions"].(map[string]interface{})
	keys := make([]*semver.Version, 0, len(versions))
	for key := range versions {
		keys = append(keys, semver.MustParse(key))
	}

	return keys
}

func FetchDependency(name string, version string, list *DependencyList) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, GetVersion(name, version))
	fmt.Println(url)
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

package graph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/harshalslimaye/ivar/internal/vercon"
)

type Dependency struct {
	Name         string
	Version      string
	Dependencies []*Dependency
}

type Version struct {
	Number     string
	Count      int
	Dependents []*Dependency
}

type Graph struct {
	Dependencies []*Dependency
	List         map[string][]*Version
}

func (g *Graph) AddToList(name string, version string, dep *Dependency) {
	_, ok := g.List[name]
	if ok {
		exists := slices.ContainsFunc(g.List[name], func(ver *Version) bool {
			return ver.Number == version
		})
		if !exists {
			g.List[name] = append(g.List[name], NewVersion(version, dep))
		} else {
			idx := VersionByIndex(g.List[name], version)
			g.List[name][idx].Count++

			haveDep := slices.ContainsFunc(g.List[name][idx].Dependents, func(d *Dependency) bool {
				return d.Name == dep.Name && d.Version == dep.Version
			})
			if !haveDep {
				g.List[name][idx].Dependents = append(g.List[name][idx].Dependents, dep)
			}
		}
	} else {
		g.List[name] = []*Version{NewVersion(version, dep)}
	}
}

func VersionByIndex(versions []*Version, version string) int {
	for i, v := range versions {
		if v.Number == version {
			return i
		}
	}

	return -1
}

func NewVersion(version string, dep *Dependency) *Version {
	return &Version{Number: version, Count: 1, Dependents: []*Dependency{dep}}
}

func BuildGraph(gph *Graph, deps map[string]string) {
	vc := vercon.NewVercon()
	for name, version := range deps {
		dep := &Dependency{Name: name, Version: vc.GetVersion(name, version), Dependencies: []*Dependency{}}
		gph.Dependencies = append(gph.Dependencies, dep)
		gph.AddToList(dep.Name, dep.Version, dep)
		FetchDependency(dep, gph, vc)
	}
}

func FetchDependency(dependency *Dependency, gph *Graph, vc *vercon.Vercon) error {
	exactVersion := vc.GetVersion(dependency.Name, dependency.Version)
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", dependency.Name, exactVersion)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch dependency %s@%s: %v", dependency.Name, dependency.Version, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch dependency %s@%s: %s", dependency.Name, dependency.Version, res.Status)
	}

	var dep struct {
		Name         string
		Version      string
		Dependencies map[string]string
	}
	err = json.NewDecoder(res.Body).Decode(&dep)
	if err != nil {
		return fmt.Errorf("failed to decode response body for dependency %s@%s: %v", dependency.Name, dependency.Version, err)
	}

	if len(dep.Dependencies) > 0 {
		for k, v := range dep.Dependencies {
			d := &Dependency{Name: k, Version: vc.GetVersion(k, v)}
			gph.AddToList(d.Name, d.Version, dependency)
			dependency.Dependencies = append(dependency.Dependencies, d)
			FetchDependency(d, gph, vc)
		}
	}

	return nil
}

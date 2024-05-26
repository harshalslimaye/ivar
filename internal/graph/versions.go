package graph

import "sync"

type Versions struct {
	Packages sync.Map
}

func NewVersions() *Versions {
	return &Versions{}
}

func (v *Versions) Set(name, version string) {
	if !v.Contains(name, version) {
		v.Packages.Store(name, append(v.Get(name), version))
	}
}

func (v *Versions) Get(name string) []string {
	if versions, exists := v.Packages.Load(name); exists {
		if values, okay := versions.([]string); okay {
			return values
		}
	}

	return []string{}
}

func (v *Versions) Contains(name, version string) bool {
	for _, v := range v.Get(name) {
		if v == version {
			return true
		}
	}

	return false
}

func (v *Versions) Len(name string) int {
	return len(v.Get(name))
}

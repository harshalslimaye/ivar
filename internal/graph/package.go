package graph

import (
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Package struct {
	Name    string
	Version string
}

func NewPackage(packageName, packageVersion string) *Package {
	return &Package{
		Name:    packageName,
		Version: registry.GetVersion(packageName, packageVersion),
	}
}

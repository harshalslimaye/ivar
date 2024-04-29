package graph

import "github.com/harshalslimaye/ivar/internal/vercon"

type Package struct {
	Name    string
	Version string
}

func NewPackage(packageName, packageVersion string) *Package {
	return &Package{
		Name:    packageName,
		Version: vercon.GetVersion(packageName, packageVersion),
	}
}

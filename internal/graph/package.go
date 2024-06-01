package graph

import (
	"fmt"

	"github.com/harshalslimaye/ivar/internal/locker"
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Package struct {
	Name       string
	Version    string
	RawVersion string
}

func NewPackage(packageName, packageVersion string, lockFile *locker.File) *Package {
	lockedVersion := lockFile.GetVersion(packageName, packageVersion)

	if lockedVersion != "" {
		return &Package{Name: packageName, Version: lockedVersion, RawVersion: packageVersion}
	}

	return &Package{
		Name:       packageName,
		Version:    registry.GetVersion(packageName, packageVersion),
		RawVersion: packageVersion,
	}
}

func (p *Package) NameAndVersion() string {
	return fmt.Sprintf("%s@%s", p.Name, p.Version)
}

package graph

import (
	"fmt"

	"github.com/harshalslimaye/ivar/internal/registry"
)

type Node struct {
	Package      *Package
	Dependencies map[string]*Node
}

func NewNode(pkg *Package) *Node {
	return &Node{
		Package:      pkg,
		Dependencies: make(map[string]*Node),
	}
}

func (n *Node) AddDependencies(deps map[string]string) {
	for depName, depVersion := range deps {
		pkg := NewPackage(depName, depVersion)
		node := NewNode(pkg)

		n.AddDependency(node)

		dependencies, err := registry.FetchDependencies(pkg.Name, pkg.Version)
		if err != nil {
			fmt.Println(err)
		}
		if len(dependencies) > 0 {
			node.AddDependencies(dependencies)
		}
	}
}

func (n *Node) AddDependency(node *Node) {
	n.Dependencies[node.Package.Name] = node
}

func (n *Node) RemoveDependency(packageName string) {
	delete(n.Dependencies, packageName)
}

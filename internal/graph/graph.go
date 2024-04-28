package graph

import (
	"fmt"

	"github.com/harshalslimaye/ivar/internal/registry"
	"github.com/harshalslimaye/ivar/internal/vercon"
)

var vc vercon.Vercon = *vercon.NewVercon()

type Package struct {
	Name    string
	Version string
}

type Node struct {
	Package      *Package
	Dependencies map[string]*Node
}

type Graph struct {
	Nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func NewDependencyGraph(deps map[string]string) *Graph {
	gh := NewGraph()

	for name, version := range deps {
		pkg := NewPackage(name, version)
		gh.AddDependencies(pkg)
	}

	return gh
}

func NewPackage(packageName, packageVersion string) *Package {
	return &Package{
		Name:    packageName,
		Version: vc.GetVersion(packageName, packageVersion),
	}
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

func (g *Graph) AddDependencies(pkg *Package) *Node {
	if node, exists := g.Nodes[pkg.Name]; exists {
		return node
	}

	node := g.AddNode(pkg)

	dependencies, err := registry.FetchDependencies(pkg.Name, pkg.Version)
	if err != nil {
		fmt.Println(err)
	}

	if len(dependencies) > 0 {
		node.AddDependencies(dependencies)
	}

	return node
}

func (n *Node) AddDependency(node *Node) {
	n.Dependencies[node.Package.Name] = node
}

func (n *Node) RemoveDependency(packageName string) {
	delete(n.Dependencies, packageName)
}

func (g *Graph) AddNode(pkg *Package) *Node {
	if _, exists := g.Nodes[pkg.Name]; !exists {
		g.Nodes[pkg.Name] = &Node{
			Package:      pkg,
			Dependencies: make(map[string]*Node),
		}
	}

	return g.Nodes[pkg.Name]
}

func (g *Graph) GetDependency(packageName string) *Package {
	if _, exists := g.Nodes[packageName]; exists {
		return g.Nodes[packageName].Package
	}

	return nil
}

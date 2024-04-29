package graph

import (
	"fmt"
	"sync"

	"github.com/harshalslimaye/ivar/internal/registry"
)

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
	var wg sync.WaitGroup

	for name, version := range deps {
		wg.Add(1)
		go func(n, v string) {
			defer wg.Done()
			pkg := NewPackage(n, v)
			gh.AddDependencies(pkg)
		}(name, version)
	}

	wg.Wait()

	return gh
}

func (g *Graph) AddDependencies(pkg *Package) *Node {
	if node, exists := g.Nodes[pkg.Name]; exists {
		return node
	}

	node := g.AddNode(pkg)

	dep, err := registry.FetchDependencies(pkg.Name, pkg.Version)
	if err != nil {
		fmt.Println(err)
	}

	node.SetBin(dep.Bin)
	if len(dep.Dependencies) > 0 {
		node.AddDependencies(dep.Dependencies)
	}

	return node
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

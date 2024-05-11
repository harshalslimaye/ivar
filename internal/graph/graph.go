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
	var mt sync.Mutex

	for name, version := range deps {
		wg.Add(1)
		go func(n, v string) {
			defer wg.Done()
			pkg := NewPackage(n, v)
			mt.Lock()
			gh.AddDependencies(pkg)
			mt.Unlock()
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

	parser, err := registry.FetchDependencies(pkg.Name, pkg.Version)
	if err != nil {
		fmt.Println(err)
	}

	node.SetMetadata(parser)
	if parser.HasDependencies() {
		node.AddDependencies(parser.GetDependencies())
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

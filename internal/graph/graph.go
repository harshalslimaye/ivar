package graph

import (
	"fmt"
	"sync"

	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Graph struct {
	Nodes     map[string]*Node
	RootNodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func NewDependencyGraph(parser *jsonparser.JsonParser) *Graph {
	gh := NewGraph()
	var wg sync.WaitGroup
	var mt sync.Mutex

	for name, version := range parser.GetObject("dependencies") {
		wg.Add(1)
		go func(n, v string) {
			defer wg.Done()
			mt.Lock()
			gh.AddDependencies(NewPackage(n, v))
			mt.Unlock()
		}(name, version)
	}

	wg.Wait()

	return gh
}

func (g *Graph) AddDependencies(pkg *Package) *Node {
	node := NewNode(pkg)

	parser, err := registry.FetchDependencies(pkg.Name, pkg.Version)
	if err != nil {
		fmt.Println(err)
	}

	node.SetMetadata(parser)
	if parser.Exists("dependencies") {
		node.AddDependencies(parser.GetObject("dependencies"))
	}

	return node
}

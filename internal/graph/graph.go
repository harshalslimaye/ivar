package graph

import (
	"fmt"
	"sync"

	"github.com/harshalslimaye/ivar/internal/constants"
	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/loader"
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Graph struct {
	Nodes            map[string]*Node
	Cache            *Cache
	Versions         *Versions
	RootDependencies []*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Cache:    NewCache(),
		Versions: NewVersions(),
	}
}

func NewDependencyGraph(parser *jsonparser.JsonParser) *Graph {
	gh := NewGraph()
	var wg sync.WaitGroup
	var mt sync.Mutex

	for _, dType := range append(constants.DEPENDENCY_TYPES, "devDependencies") {
		for name, version := range parser.GetObject(dType) {
			wg.Add(1)
			go func(n, v, t string) {
				defer wg.Done()
				mt.Lock()
				gh.AddDependencies(NewPackage(n, v), t)
				mt.Unlock()
			}(name, version, dType)
		}
	}

	wg.Wait()
	loader.Clear()

	return gh
}

func (g *Graph) AddDependencies(pkg *Package, category string) {
	node := NewNode(pkg, category, g)
	g.AtRoot(node)

	parser, err := registry.FetchDependencies(pkg.Name, pkg.Version)
	if err != nil {
		fmt.Println(err)
	} else {
		g.Nodes[pkg.Name] = node
		node.SetMetadata(parser)
		for _, dType := range constants.DEPENDENCY_TYPES {
			if parser.Exists(dType) {
				node.AddDependencies(parser.GetObject("dependencies"), dType)
			}
		}
	}
}

func (g *Graph) AtRoot(n *Node) {
	g.RootDependencies = append(g.RootDependencies, n)
}

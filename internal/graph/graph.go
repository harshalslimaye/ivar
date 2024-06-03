package graph

import (
	"fmt"

	"github.com/alitto/pond"
	"github.com/harshalslimaye/ivar/internal/cache"
	"github.com/harshalslimaye/ivar/internal/constants"
	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/loader"
	"github.com/harshalslimaye/ivar/internal/locker"
	"github.com/harshalslimaye/ivar/internal/registry"
	"github.com/harshalslimaye/ivar/internal/store"
)

type Graph struct {
	Nodes            map[string]*Node
	Store            *store.Store
	Versions         *Versions
	RootDependencies []*Node
	LockFile         *locker.File
	HasCache         bool
	Cache            *cache.Cache
	Pool             *pond.WorkerPool
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Store:    store.GetStore(),
		Versions: NewVersions(),
		Cache:    cache.NewCache(),
		LockFile: locker.NewLocker(),
		Pool:     pond.New(25, 0, pond.MinWorkers(10)),
	}
}

func NewDependencyGraph(parser *jsonparser.JsonParser) *Graph {
	gh := NewGraph()

	for _, dType := range append(constants.DEPENDENCY_TYPES, "devDependencies") {
		for name, version := range parser.GetObject(dType) {

			gh.Pool.Submit(func() {
				gh.AddDependencies(NewPackage(name, version, gh.LockFile), dType)
			})
		}
	}

	loader.Clear()

	gh.Pool.StopAndWait()

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

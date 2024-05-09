package graph

import (
	"fmt"
	"sync"

	"github.com/harshalslimaye/ivar/internal/registry"
)

type Node struct {
	Package      *Package
	Dependencies map[string]*Node
	Bin          map[string]string
	mutex        sync.Mutex
}

func NewNode(pkg *Package) *Node {
	return &Node{
		Package:      pkg,
		Dependencies: make(map[string]*Node),
		Bin:          make(map[string]string),
	}
}

func (n *Node) AddDependencies(deps map[string]string) {
	var wg sync.WaitGroup

	for depName, depVersion := range deps {
		wg.Add(1)

		go func(name, version string) {
			defer wg.Done()

			pkg := NewPackage(name, version)
			node := NewNode(pkg)

			n.Lock()
			n.AddDependency(node)
			n.Unlock()

			dep, err := registry.FetchDependencies(pkg.Name, pkg.Version)
			if err != nil {
				fmt.Println(err)
			}

			node.SetBin(dep.Bin)
			if len(dep.Dependencies) > 0 {
				node.AddDependencies(dep.Dependencies)
			}
		}(depName, depVersion)
	}

	wg.Wait()
}

func (n *Node) AddDependency(node *Node) {
	n.Dependencies[node.Package.Name] = node
}

func (n *Node) SetBin(bin map[string]string) {
	if len(bin) > 0 {
		n.Bin = bin
	}
}

func (n *Node) Lock() {
	n.mutex.Lock()
}

func (n *Node) Unlock() {
	n.mutex.Unlock()
}

package graph

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"sync"

	"github.com/harshalslimaye/ivar/internal/constants"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/loader"
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Node struct {
	Graph        *Graph
	Parent       *Node
	Package      *Package
	Dependencies map[string]*Node
	Bin          map[string]string
	TarballUrl   string
	FileName     string
	Integrity    string
	Category     string
	mutex        sync.Mutex
	PrunedNodes  map[string]*Node
}

func NewNode(pkg *Package, category string, gh *Graph) *Node {
	node := gh.Store.Get(pkg.NameAndVersion())
	gh.Versions.Set(pkg.Name, pkg.Version)

	if node == nil {
		node = &Node{
			Package:      pkg,
			Dependencies: make(map[string]*Node),
			Bin:          make(map[string]string),
			Category:     category,
			Graph:        gh,
			PrunedNodes:  make(map[string]*Node),
		}
		gh.Store.Set(pkg.NameAndVersion(), node)
	}

	return node
}

func (n *Node) AddDependencies(deps map[string]string, category string) {
	loader.Show("\r" + "Resolving " + n.Name() + "@" + n.Version() + "...")
	var wg sync.WaitGroup

	for depName, depVersion := range deps {
		wg.Add(1)

		go func(name, version string) {
			defer wg.Done()
			n.Lock()
			node := NewNode(NewPackage(name, version, n.Graph.LockFile), category, n.Graph)
			n.Unlock()

			parser, err := registry.FetchDependencies(node.Name(), node.Version())
			if err != nil {
				fmt.Printf("failed to resolve %s@%s: %s \n", node.Name(), node.Version(), err.Error())
				fmt.Println(err)
				return
			}

			n.Lock()
			n.AddDependency(node)
			if category == "peerDependencies" {
				n.Graph.AtRoot(node)
			}
			n.Unlock()

			node.SetMetadata(parser)
			for _, dType := range constants.DEPENDENCY_TYPES {
				if parser.Exists(dType) {
					node.AddDependencies(parser.GetObject(dType), dType)
				}
			}

		}(depName, depVersion)
	}

	wg.Wait()
}

func (n *Node) AddDependency(node *Node) {
	n.Dependencies[node.Package.Name] = node
}

func (n *Node) SetMetadata(parser *jsonparser.JsonParser) {
	n.Lock()
	n.TarballUrl = parser.TarballUrl()
	n.Integrity = parser.GetValue("integrity")
	n.SetTarName(n.TarballUrl)
	n.SetBin(parser.GetBin())
	n.Unlock()
}

func (n *Node) SetTarName(urlString string) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	n.FileName = path.Base(parsedURL.Path)
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

func (n *Node) Name() string {
	return n.Package.Name
}

func (n *Node) Version() string {
	return n.Package.Version
}

func (n *Node) SourcePath(dir string) string {
	return filepath.Join(dir, n.FileName)
}

func (n *Node) TargetPath(dir string) string {
	return filepath.Join(helper.GetCurrentDirPath(), dir)
}

func (n *Node) IsPeer() bool {
	return n.Category == "peerDependencies"
}

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
	DownloadDir  string
	Integrity    string
	Category     string
	mutex        sync.Mutex
}

func NewNode(pkg *Package, category string) *Node {
	return &Node{
		Package:      pkg,
		Dependencies: make(map[string]*Node),
		Bin:          make(map[string]string),
		Category:     category,
	}
}

func (n *Node) AddDependencies(deps map[string]string, category string) {
	var wg sync.WaitGroup

	for depName, depVersion := range deps {
		wg.Add(1)

		go func(name, version string) {
			defer wg.Done()

			pkg := NewPackage(name, version)
			node := NewNode(pkg, category)

			n.Lock()
			n.AddDependency(node)
			n.Unlock()

			parser, err := registry.FetchDependencies(pkg.Name, pkg.Version)
			if err != nil {
				fmt.Printf("Failed to download %s@%s: \n", pkg.Name, pkg.Version)
				fmt.Println(err)
			}

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
	n.TarballUrl = parser.TarballUrl()
	n.Integrity = parser.GetValue("integrity")
	n.SetTarName(n.TarballUrl)
	n.SetBin(parser.GetBin())
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

func (n *Node) SourcePath() string {
	return filepath.Join(n.DownloadDir, n.FileName)
}

func (n *Node) TargetPath() string {
	return filepath.Join(helper.GetCurrentDirPath(), n.DownloadDir)
}

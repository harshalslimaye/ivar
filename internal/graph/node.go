package graph

import (
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/registry"
)

type Node struct {
	Package      *Package
	Dependencies map[string]*Node
	Bin          map[string]string
	DownloadPath string
	FileName     string
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

			parser, err := registry.FetchDependencies(pkg.Name, pkg.Version)
			if err != nil {
				fmt.Printf("Failed to download %s@%s: \n", pkg.Name, pkg.Version)
				fmt.Println(err)
			}

			node.SetMetadata(parser)
			if len(parser.GetDependencies()) > 0 {
				node.AddDependencies(parser.GetDependencies())
			}
		}(depName, depVersion)
	}

	wg.Wait()
}

func (n *Node) AddDependency(node *Node) {
	n.Dependencies[node.Package.Name] = node
}

func (n *Node) SetMetadata(parser *jsonparser.JsonParser) {
	n.DownloadPath = parser.GetDownloadPath()
	n.SetTarName(n.DownloadPath)
	n.SetBin(parser.GetBin())
}

func (n *Node) SetTarName(urlString string) {
	// Parse the URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get the base name of the path
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

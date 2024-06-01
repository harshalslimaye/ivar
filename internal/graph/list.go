package graph

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Masterminds/semver/v3"
)

type List struct {
	NodeMap sync.Map
}

func (l *List) Map() map[string]*Node {
	nodes := make(map[string]*Node)

	l.NodeMap.Range(func(key, value interface{}) bool {
		node, isNode := value.(*Node)
		path, isString := key.(string)

		if isNode && isString {
			nodes[path] = node
		}
		return true
	})

	return nodes
}

func (l *List) AddNode(n *Node, p *Node) {
	rootPath := filepath.Join("node_modules", n.Name())
	value, exists := l.NodeMap.Load(rootPath)

	if exists {
		if node, okay := value.(*Node); !okay || node.Version() != n.Version() {
			p.PrunedNodes[filepath.Join("node_modules", p.Name(), "node_modules", n.Name())] = n
			return
		}
	}

	if n.Graph.Versions.Len(n.Name()) > 1 {
		var versions []*semver.Version
		for _, v := range n.Graph.Versions.Get(n.Name()) {
			version, err := semver.NewVersion(v)
			if err != nil {
				fmt.Printf("Error parsing version %s: %v\n", v, err)
				continue
			}
			versions = append(versions, version)
		}

		sort.Slice(versions, func(i, j int) bool {
			return versions[i].LessThan(versions[j])
		})

		if len(versions) > 0 {
			latestVersion := versions[len(versions)-1]
			if latestVersion.String() == n.Version() {
				l.NodeMap.Store(rootPath, n)
			} else {
				p.PrunedNodes[filepath.Join("node_modules", p.Name(), "node_modules", n.Name())] = n
			}
			return
		}
	}

	if !exists {
		l.NodeMap.Store(rootPath, n)
		return
	}
}

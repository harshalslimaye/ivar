package graph

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/harshalslimaye/ivar/internal/locker"
)

type List struct {
	NodeMap sync.Map
	Graph   *Graph
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
			key := filepath.Join("node_modules", p.Name(), "node_modules", n.Name())
			p.PrunedNodes[key] = n
			l.Graph.LockFile.Add(NewLockItem(n, key), key)
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
				l.Graph.LockFile.Add(NewLockItem(n, rootPath), rootPath)
			} else {
				key := filepath.Join("node_modules", p.Name(), "node_modules", n.Name())
				p.PrunedNodes[key] = n
				l.Graph.LockFile.Add(NewLockItem(n, key), key)
			}
			return
		}
	}

	if !exists {
		l.NodeMap.Store(rootPath, n)
		l.Graph.LockFile.Add(NewLockItem(n, rootPath), rootPath)
		return
	}
}

func NewList(gh *Graph) *List {
	return &List{
		NodeMap: sync.Map{},
		Graph:   gh,
	}
}

func NewLockItem(node *Node, path string) *locker.Element {
	return &locker.Element{
		Version:   node.Version(),
		Resolved:  node.TarballUrl,
		Integrity: node.Integrity,
		Path:      path,
	}
}

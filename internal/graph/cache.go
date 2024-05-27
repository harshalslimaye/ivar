package graph

import "sync"

type Store struct {
	Nodes sync.Map
}

func (c *Store) Get(key string) *Node {
	if node, exists := c.Nodes.Load(key); exists {
		if n, ok := node.(*Node); ok {
			return n
		}
	}

	return nil
}

func (c *Store) Set(key string, value *Node) {
	c.Nodes.Store(key, value)
}

func NewStore() *Store {
	return &Store{}
}

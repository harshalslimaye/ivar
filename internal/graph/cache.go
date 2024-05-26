package graph

import "sync"

type Cache struct {
	Nodes sync.Map
}

func (c *Cache) Get(key string) *Node {
	if node, exists := c.Nodes.Load(key); exists {
		if n, ok := node.(*Node); ok {
			return n
		}
	}

	return nil
}

func (c *Cache) Set(key string, value *Node) {
	c.Nodes.Store(key, value)
}

func NewCache() *Cache {
	return &Cache{}
}

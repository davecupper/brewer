package graph

import (
	"fmt"

	"github.com/nickcorin/brewer/internal/config"
)

// Graph represents a directed acyclic graph of service dependencies.
type Graph struct {
	nodes map[string]*Node
}

// Node represents a single service in the dependency graph.
type Node struct {
	Service  config.Service
	Deps     []*Node
	visited  bool
	inStack  bool
}

// Build constructs a dependency graph from the provided config.
func Build(cfg *config.Config) (*Graph, error) {
	g := &Graph{
		nodes: make(map[string]*Node, len(cfg.Services)),
	}

	for _, svc := range cfg.Services {
		g.nodes[svc.Name] = &Node{Service: svc}
	}

	for _, svc := range cfg.Services {
		node := g.nodes[svc.Name]
		for _, dep := range svc.DependsOn {
			depNode, ok := g.nodes[dep]
			if !ok {
				return nil, fmt.Errorf("service %q depends on unknown service %q", svc.Name, dep)
			}
			node.Deps = append(node.Deps, depNode)
		}
	}

	if err := g.detectCycles(); err != nil {
		return nil, err
	}

	return g, nil
}

// TopologicalOrder returns services in the order they should be started.
func (g *Graph) TopologicalOrder() ([]config.Service, error) {
	var order []config.Service
	visited := make(map[string]bool)

	var visit func(n *Node)
	visit = func(n *Node) {
		if visited[n.Service.Name] {
			return
		}
		visited[n.Service.Name] = true
		for _, dep := range n.Deps {
			visit(dep)
		}
		order = append(order, n.Service)
	}

	for _, node := range g.nodes {
		visit(node)
	}

	return order, nil
}

func (g *Graph) detectCycles() error {
	for _, node := range g.nodes {
		node.visited = false
		node.inStack = false
	}
	for _, node := range g.nodes {
		if !node.visited {
			if err := dfs(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func dfs(n *Node) error {
	n.visited = true
	n.inStack = true
	for _, dep := range n.Deps {
		if !dep.visited {
			if err := dfs(dep); err != nil {
				return err
			}
		} else if dep.inStack {
			return fmt.Errorf("cycle detected involving service %q", dep.Service.Name)
		}
	}
	n.inStack = false
	return nil
}

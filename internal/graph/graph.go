// internal/graph/graph.go
package graph

type Graph struct {
    Nodes    map[string]bool
    Edges    map[string]map[string]bool
    Start    string
    End      string
}

func NewGraph() *Graph {
    return &Graph{
        Nodes: make(map[string]bool),
        Edges: make(map[string]map[string]bool),
    }
}

func (g *Graph) AddNode(name string) {
    if !g.Nodes[name] {
        g.Nodes[name] = true
        g.Edges[name] = make(map[string]bool)
    }
}

func (g *Graph) AddEdge(from, to string) {
    g.AddNode(from)
    g.AddNode(to)
    g.Edges[from][to] = true
    g.Edges[to][from] = true
}
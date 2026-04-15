// internal/graph/edmonds_karp.go
package graph

type FlowEdge struct {
    To     string
    Rev    int  // Index of reverse edge in adjacency list
    Cap    int  // Capacity remaining
}

type FlowGraph struct {
    Adj    map[string][]FlowEdge
    Nodes  map[string]bool
}

func NewFlowGraph() *FlowGraph {
    return &FlowGraph{
        Adj:   make(map[string][]FlowEdge),
        Nodes: make(map[string]bool),
    }
}

func (fg *FlowGraph) AddNode(name string) {
    if !fg.Nodes[name] {
        fg.Nodes[name] = true
        fg.Adj[name] = []FlowEdge{}
    }
}

func (fg *FlowGraph) AddEdge(from, to string, cap int) {
    fg.AddNode(from)
    fg.AddNode(to)
    
    forward := FlowEdge{To: to, Rev: len(fg.Adj[to]), Cap: cap}
    backward := FlowEdge{To: from, Rev: len(fg.Adj[from]), Cap: 0}
    
    fg.Adj[from] = append(fg.Adj[from], forward)
    fg.Adj[to] = append(fg.Adj[to], backward)
}

func (fg *FlowGraph) MaxFlow(source, sink string) (int, [][]string) {
    flow := 0
    paths := [][]string{}
    
    for {
        parent := make(map[string]*FlowEdge)
        found := fg.bfs(source, sink, parent)
        
        if !found {
            break
        }
        
        // Find the path
        path := []string{}
        for v := sink; v != source; v = parent[v].To {
            path = append([]string{v}, path...)
        }
        path = append([]string{source}, path...)
        paths = append(paths, path)
        
        // Update capacities (minimum capacity is 1 for our graph)
        for v := sink; v != source; v = parent[v].To {
            edge := parent[v]
            edge.Cap--
            // Update reverse edge
            rev := &fg.Adj[edge.To][edge.Rev]
            rev.Cap++
        }
        flow++
    }
    
    return flow, paths
}

func (fg *FlowGraph) bfs(source, sink string, parent map[string]*FlowEdge) bool {
    visited := make(map[string]bool)
    queue := []string{source}
    visited[source] = true
    
    for len(queue) > 0 {
        u := queue[0]
        queue = queue[1:]
        
        for i := range fg.Adj[u] {
            edge := &fg.Adj[u][i]
            if !visited[edge.To] && edge.Cap > 0 {
                visited[edge.To] = true
                parent[edge.To] = edge
                if edge.To == sink {
                    return true
                }
                queue = append(queue, edge.To)
            }
        }
    }
    return false
}
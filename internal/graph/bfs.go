// internal/graph/bfs.go
package graph

type Path struct {
    Nodes []string
}

func (g *Graph) FindShortestPath() *Path {
    if g.Start == "" || g.End == "" {
        return nil
    }

    queue := [][]string{{g.Start}}
    visited := map[string]bool{g.Start: true}

    for len(queue) > 0 {
        path := queue[0]
        queue = queue[1:]
        node := path[len(path)-1]

        if node == g.End {
            return &Path{Nodes: path}
        }

        for neighbor := range g.Edges[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                newPath := make([]string, len(path))
                copy(newPath, path)
                newPath = append(newPath, neighbor)
                queue = append(queue, newPath)
            }
        }
    }
    return nil
}
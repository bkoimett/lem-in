// internal/graph/edmonds_karp_test.go
package graph

import (
    "testing"
)

func TestMaxFlowSimple(t *testing.T) {
    fg := NewFlowGraph()
    
    // Simple graph: 0 -> 1 -> 2, capacity 1 each
    fg.AddEdge("0", "1", 1)
    fg.AddEdge("1", "2", 1)
    
    flow, paths := fg.MaxFlow("0", "2")
    
    if flow != 1 {
        t.Errorf("Expected flow 1, got %d", flow)
    }
    
    if len(paths) != 1 {
        t.Errorf("Expected 1 path, got %d", len(paths))
    }
}

func TestMaxFlowMultiplePaths(t *testing.T) {
    fg := NewFlowGraph()
    
    // Graph with two parallel paths: 0->1->3 and 0->2->3
    fg.AddEdge("0", "1", 1)
    fg.AddEdge("1", "3", 1)
    fg.AddEdge("0", "2", 1)
    fg.AddEdge("2", "3", 1)
    
    flow, paths := fg.MaxFlow("0", "3")
    
    if flow != 2 {
        t.Errorf("Expected flow 2, got %d", flow)
    }
    
    if len(paths) != 2 {
        t.Errorf("Expected 2 paths, got %d", len(paths))
    }
}
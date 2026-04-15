// internal/solver/solver_test.go
package solver

import (
	"fmt"
    "strings"
    "testing"
)

func TestSimulationSinglePath(t *testing.T) {
    paths := [][]string{
        {"start", "A", "B", "end"},
    }
    
    sim := NewSimulation(paths, 3)
    sim.Run()
    
    moves := sim.String()
    
    // Check that all ants reach end
    antPositions := make(map[int]string)
    lines := strings.Split(moves, "\n")
    for _, line := range lines {
        if line == "" {
            continue
        }
        moves := strings.Split(line, " ")
        for _, move := range moves {
            var antID int
            var room string
            fmt.Sscanf(move, "L%d-%s", &antID, &room)
            antPositions[antID] = room
        }
    }
    
    // Last move should have all ants at end
    if len(lines) > 0 {
        lastLine := lines[len(lines)-1]
        moves := strings.Split(lastLine, " ")
        for _, move := range moves {
            var antID int
            var room string
            fmt.Sscanf(move, "L%d-%s", &antID, &room)
            if room != "end" {
                t.Errorf("Ant %d not at end, at %s", antID, room)
            }
        }
    }
}

func TestSimulationMultiplePaths(t *testing.T) {
    paths := [][]string{
        {"start", "A", "end"},
        {"start", "B", "end"},
    }
    
    sim := NewSimulation(paths, 4)
    sim.Run()
    
    // Should be more efficient than single path
    moves := strings.Split(sim.String(), "\n")
    
    // With 2 parallel paths of length 2, 4 ants should take 3 turns
    if len(moves) > 4 {
        t.Errorf("Took too many turns: %d", len(moves))
    }
}
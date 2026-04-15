// internal/solver/solver.go
package solver

import (
    "fmt"
    "strings"
)

type Ant struct {
    ID    int
    Path  []string
    Index int // Current position index
}

type Simulation struct {
    Ants      []*Ant
    Occupied  map[string]int // Room -> ant ID that occupies it
    Start     string
    End       string
    Moves     [][]string // Each turn's moves: ["L1-room", "L2-room"]
}

func NewSimulation(paths [][]string, numAnts int) *Simulation {
    sim := &Simulation{
        Ants:     make([]*Ant, 0),
        Occupied: make(map[string]int),
        Moves:    make([][]string, 0),
    }
    
    // Assign ants to paths (simple round-robin)
    antID := 1
    for antID <= numAnts {
        for _, path := range paths {
            if antID > numAnts {
                break
            }
            // Don't include start room in the path for movement
            pathForAnt := path[1:] // Exclude start
            sim.Ants = append(sim.Ants, &Ant{
                ID:    antID,
                Path:  pathForAnt,
                Index: -1, // Not yet started
            })
            antID++
        }
    }
    
    if len(paths) > 0 {
        sim.Start = paths[0][0]
        sim.End = paths[0][len(paths[0])-1]
    }
    
    return sim
}

func (s *Simulation) Run() {
    antsInEnd := 0
    totalAnts := len(s.Ants)
    
    for antsInEnd < totalAnts {
        turn := make([]string, 0)
        newOccupied := make(map[string]int)
        
        // Copy occupied from end room (end can have multiple ants)
        for room, antID := range s.Occupied {
            if room == s.End {
                newOccupied[room] = antID
            }
        }
        
        // Move ants
        for _, ant := range s.Ants {
            if ant.Index == len(ant.Path)-1 {
                // Already at end
                continue
            }
            
            nextIndex := ant.Index + 1
            if nextIndex >= len(ant.Path) {
                continue
            }
            
            nextRoom := ant.Path[nextIndex]
            
            // Check if room is occupied (and not end)
            if _, occupied := newOccupied[nextRoom]; occupied && nextRoom != s.End {
                continue
            }
            
            // Move the ant
            ant.Index = nextIndex
            newOccupied[nextRoom] = ant.ID
            turn = append(turn, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
            
            if nextRoom == s.End {
                antsInEnd++
            }
        }
        
        s.Occupied = newOccupied
        if len(turn) > 0 {
            s.Moves = append(s.Moves, turn)
        }
    }
}

func (s *Simulation) String() string {
    result := make([]string, len(s.Moves))
    for i, turn := range s.Moves {
        result[i] = strings.Join(turn, " ")
    }
    return strings.Join(result, "\n")
}
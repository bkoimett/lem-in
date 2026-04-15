// Package tests contains unit tests for the solver package.
package tests

import (
	"lemin/colony"
	"lemin/solver"
	"strings"
	"testing"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

// buildColony is a test helper that creates a Colony from a simple
// adjacency description. Rooms are names a, b, c... with dummy coordinates.
func buildColony(numAnts int, startRoom, endRoom string, rooms []string, links [][2]string) *colony.Colony {
	c := colony.NewColony()
	c.NumAnts = numAnts
	c.StartRoom = startRoom
	c.EndRoom = endRoom
	for _, name := range rooms {
		c.Rooms[name] = &colony.Room{Name: name}
	}
	c.Rooms[startRoom].IsStart = true
	c.Rooms[endRoom].IsEnd = true
	for _, lk := range links {
		a, b := lk[0], lk[1]
		c.Rooms[a].Links = append(c.Rooms[a].Links, b)
		c.Rooms[b].Links = append(c.Rooms[b].Links, a)
	}
	return c
}

// countTurns returns the number of turns in the simulation output.
func countTurns(lines []string) int {
	count := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "L") {
			count++
		}
	}
	return count
}

// allAntsReachedEnd verifies that every ant (1..n) appears exactly once
// in a move to the end room name.
func allAntsReachedEnd(moves []string, numAnts int, endRoom string) bool {
	suffix := "-" + endRoom
	reached := make(map[int]bool)
	for _, line := range moves {
		for _, part := range strings.Fields(line) {
			if strings.HasSuffix(part, suffix) {
				// Parse ant ID
				var id int
				for _, ch := range part[1:] {
					if ch == '-' {
						break
					}
					id = id*10 + int(ch-'0')
				}
				reached[id] = true
			}
		}
	}
	return len(reached) == numAnts
}

// ── FindPaths Tests ───────────────────────────────────────────────────────────

func TestFindPathsSimple(t *testing.T) {
	c := buildColony(1, "S", "E",
		[]string{"S", "E"},
		[][2]string{{"S", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("expected at least one path")
	}
	if paths[0][0] != "S" || paths[0][len(paths[0])-1] != "E" {
		t.Error("path should go from S to E")
	}
}

func TestFindPathsNoPath(t *testing.T) {
	c := buildColony(2, "S", "E",
		[]string{"S", "mid", "E"},
		[][2]string{{"S", "mid"}}, // no connection to E
	)
	_, err := solver.FindPaths(c)
	if err == nil {
		t.Error("expected error when no path exists")
	}
}

func TestFindPathsMultiplePaths(t *testing.T) {
	// Diamond graph: S -> A -> E  and  S -> B -> E
	c := buildColony(2, "S", "E",
		[]string{"S", "A", "B", "E"},
		[][2]string{{"S", "A"}, {"S", "B"}, {"A", "E"}, {"B", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With 2 ants and 2 paths of equal length, both paths should be used
	if len(paths) < 2 {
		t.Errorf("expected 2 paths for diamond graph, got %d", len(paths))
	}
}

func TestFindPathsShortestFirst(t *testing.T) {
	// Two paths: S->E (length 1) and S->M1->M2->E (length 3)
	c := buildColony(1, "S", "E",
		[]string{"S", "M1", "M2", "E"},
		[][2]string{{"S", "E"}, {"S", "M1"}, {"M1", "M2"}, {"M2", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The direct path S->E should be first (or only path used)
	if len(paths[0]) != 2 {
		t.Errorf("expected shortest path first (length 2), got length %d", len(paths[0]))
	}
}

// ── SimulateAnts Tests ────────────────────────────────────────────────────────

func TestSimulateAntsAllReachEnd(t *testing.T) {
	c := buildColony(4, "S", "E",
		[]string{"S", "A", "E"},
		[][2]string{{"S", "A"}, {"A", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("FindPaths error: %v", err)
	}
	moves := solver.SimulateAnts(4, paths)
	if !allAntsReachedEnd(moves, 4, "E") {
		t.Error("not all ants reached the end room")
	}
}

func TestSimulateAntsNoCollisions(t *testing.T) {
	// Build a linear chain: S-A-B-C-E
	// With multiple ants, they should never be in the same room at the same turn
	c := buildColony(5, "S", "E",
		[]string{"S", "A", "B", "C", "E"},
		[][2]string{{"S", "A"}, {"A", "B"}, {"B", "C"}, {"C", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("FindPaths error: %v", err)
	}
	moves := solver.SimulateAnts(5, paths)

	// Track ant positions per turn
	for _, line := range moves {
		rooms := make(map[string][]int) // room -> list of ant IDs
		for _, part := range strings.Fields(line) {
			dashIdx := strings.LastIndex(part, "-")
			if dashIdx < 0 {
				continue
			}
			room := part[dashIdx+1:]
			var id int
			for _, ch := range part[1:dashIdx] {
				id = id*10 + int(ch-'0')
			}
			// Skip start and end — they allow multiple ants
			if room == "S" || room == "E" {
				continue
			}
			rooms[room] = append(rooms[room], id)
		}
		for room, ants := range rooms {
			if len(ants) > 1 {
				t.Errorf("collision in room %s on turn (line: %s): ants %v", room, line, ants)
			}
		}
	}
}

func TestExample00TurnCount(t *testing.T) {
	// example00: 4 ants, path S->2->3->1, should take at most 6 turns
	c := buildColony(4, "0", "1",
		[]string{"0", "2", "3", "1"},
		[][2]string{{"0", "2"}, {"2", "3"}, {"3", "1"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("FindPaths error: %v", err)
	}
	moves := solver.SimulateAnts(4, paths)
	turns := len(moves)
	if turns > 6 {
		t.Errorf("example00: expected <=6 turns, got %d", turns)
	}
}

func TestSingleAntDirectPath(t *testing.T) {
	c := buildColony(1, "S", "E",
		[]string{"S", "E"},
		[][2]string{{"S", "E"}},
	)
	paths, _ := solver.FindPaths(c)
	moves := solver.SimulateAnts(1, paths)
	if len(moves) != 1 {
		t.Errorf("single ant direct path should take exactly 1 turn, got %d", len(moves))
	}
}

func TestCalculateTurnsDiamondGraph(t *testing.T) {
	// Diamond: two paths each length 2, 2 ants -> should take 2 turns
	c := buildColony(2, "S", "E",
		[]string{"S", "A", "B", "E"},
		[][2]string{{"S", "A"}, {"S", "B"}, {"A", "E"}, {"B", "E"}},
	)
	paths, err := solver.FindPaths(c)
	if err != nil {
		t.Fatalf("FindPaths error: %v", err)
	}
	moves := solver.SimulateAnts(2, paths)
	if len(moves) > 2 {
		t.Errorf("diamond with 2 ants: expected <=2 turns, got %d", len(moves))
	}
}

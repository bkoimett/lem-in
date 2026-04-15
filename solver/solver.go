// Package solver finds optimal ant paths and simulates movement.
//
// # The Core Problem
//
// We need to send N ants from ##start to ##end through a graph where:
//   - Each interior room can hold only ONE ant at a time.
//   - Each tunnel can be used only once per turn.
//   - We want to minimise the total number of turns.
//
// # Approach: Node-Split Max-Flow (Edmonds-Karp)
//
// To enforce the "one ant per room" rule, we use a standard graph theory trick:
// split every interior node v into two nodes: v_in and v_out, connected by a
// directed edge with capacity 1. All edges coming INTO v connect to v_in, and
// all edges going OUT of v come from v_out.
//
// This transforms the room-capacity problem into a standard edge-capacity
// max-flow problem, which Edmonds-Karp (BFS-based Ford-Fulkerson) solves.
//
// After finding max-flow, we extract the actual paths by reading the flow
// assignments, then simulate ant movement optimally.
package solver

import (
	"fmt"
	"lemin/colony"
)

// FindPaths finds the optimal set of paths from start to end.
// It returns the paths sorted shortest-first.
func FindPaths(c *colony.Colony) ([][]string, error) {
	// Build the flow network from the colony graph
	net := buildFlowNetwork(c)

	// Run Edmonds-Karp to find max flow (= max number of parallel paths)
	runMaxFlow(net, c.StartRoom, c.EndRoom)

	// Extract individual paths from the flow assignments
	paths, err := extractPaths(net, c, c.StartRoom, c.EndRoom)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, no path between start and end")
	}

	// Sort paths shortest first (bubble sort — paths are few)
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if len(paths[j]) < len(paths[i]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	// Select the best subset given our ant count
	best := selectBestPaths(paths, c.NumAnts)
	return best, nil
}

// ── Flow Network ─────────────────────────────────────────────────────────────

// edge represents a directed edge in the flow network.
// Each undirected tunnel becomes two directed edges (forward + backward).
// Each node-split creates one internal edge (v_in → v_out).
type edge struct {
	to  int // destination node index
	cap int // capacity
	flow int // current flow
	rev int // index of the reverse edge in the adjacency list of `to`
}

// flowNetwork is the graph used for max-flow computation.
type flowNetwork struct {
	adj      [][]edge
	nodeID   map[string]int
	numNodes int
}

// nodeIn returns the "entry" node index for a room.
func (fn *flowNetwork) nodeIn(name string) int {
	return fn.nodeID[name+"_in"]
}

// nodeOut returns the "exit" node index for a room.
func (fn *flowNetwork) nodeOut(name string) int {
	return fn.nodeID[name+"_out"]
}

// addEdge adds a directed edge u->v with given capacity, plus its reverse edge.
func (fn *flowNetwork) addEdge(u, v, cap int) {
	fn.adj[u] = append(fn.adj[u], edge{to: v, cap: cap, rev: len(fn.adj[v])})
	fn.adj[v] = append(fn.adj[v], edge{to: u, cap: 0, rev: len(fn.adj[u]) - 1})
}

// buildFlowNetwork constructs the node-split flow graph from the colony.
//
// For every room R:
//   - Create R_in  (receives incoming flow)
//   - Create R_out (sends outgoing flow)
//   - Add edge R_in -> R_out with capacity 1
//     (capacity large for ##start and ##end so they don't limit flow)
//
// For every tunnel A<->B:
//   - Add edge A_out -> B_in  with capacity 1
//   - Add edge B_out -> A_in  with capacity 1  (undirected)
func buildFlowNetwork(c *colony.Colony) *flowNetwork {
	nodeID := make(map[string]int)
	idx := 0
	for name := range c.Rooms {
		nodeID[name+"_in"] = idx
		idx++
		nodeID[name+"_out"] = idx
		idx++
	}

	fn := &flowNetwork{
		adj:      make([][]edge, idx),
		nodeID:   nodeID,
		numNodes: idx,
	}

	// Internal edges: R_in -> R_out
	for name := range c.Rooms {
		inNode := fn.nodeIn(name)
		outNode := fn.nodeOut(name)
		cap := 1
		if name == c.StartRoom || name == c.EndRoom {
			cap = len(c.Rooms) + c.NumAnts
		}
		fn.addEdge(inNode, outNode, cap)
	}

	// Tunnel edges
	seen := make(map[string]bool)
	for nameA, room := range c.Rooms {
		for _, nameB := range room.Links {
			key := nameA + "|" + nameB
			if nameB < nameA {
				key = nameB + "|" + nameA
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			fn.addEdge(fn.nodeOut(nameA), fn.nodeIn(nameB), 1)
			fn.addEdge(fn.nodeOut(nameB), fn.nodeIn(nameA), 1)
		}
	}

	return fn
}

// ── Edmonds-Karp (BFS Ford-Fulkerson) ────────────────────────────────────────

// runMaxFlow runs Edmonds-Karp max-flow from start_out to end_in.
// It repeatedly finds shortest augmenting paths via BFS and pushes flow.
func runMaxFlow(fn *flowNetwork, startName, endName string) {
	source := fn.nodeOut(startName)
	sink := fn.nodeIn(endName)

	for {
		parent := make([]int, fn.numNodes)
		parentEdge := make([]int, fn.numNodes)
		for i := range parent {
			parent[i] = -1
		}
		parent[source] = source
		queue := []int{source}

		for len(queue) > 0 && parent[sink] == -1 {
			u := queue[0]
			queue = queue[1:]
			for i, e := range fn.adj[u] {
				if parent[e.to] == -1 && e.cap-e.flow > 0 {
					parent[e.to] = u
					parentEdge[e.to] = i
					queue = append(queue, e.to)
				}
			}
		}

		if parent[sink] == -1 {
			break // max flow reached
		}

		// Push 1 unit of flow along the found path
		v := sink
		for v != source {
			u := parent[v]
			ei := parentEdge[v]
			fn.adj[u][ei].flow++
			fn.adj[v][fn.adj[u][ei].rev].flow--
			v = u
		}
	}
}

// ── Path Extraction ───────────────────────────────────────────────────────────

// extractPaths reads the flow assignments and reconstructs the actual paths.
// Each unit of flow through start_out represents one path.
func extractPaths(fn *flowNetwork, c *colony.Colony, startName, endName string) ([][]string, error) {
	// Reverse lookup: node index -> room name (for _in nodes)
	roomOfIn := make(map[int]string)
	for name := range c.Rooms {
		roomOfIn[fn.nodeIn(name)] = name
	}

	var paths [][]string
	startOut := fn.nodeOut(startName)

	for i, e := range fn.adj[startOut] {
		if e.flow <= 0 {
			continue
		}

		// Consume this flow unit
		fn.adj[startOut][i].flow--
		fn.adj[e.to][e.rev].flow++

		path := []string{startName}
		cur := e.to // a room_in node

		for {
			roomName, ok := roomOfIn[cur]
			if !ok {
				return nil, fmt.Errorf("ERROR: flow extraction failed — unknown node")
			}
			path = append(path, roomName)
			if roomName == endName {
				break
			}

			// Move from roomName_in to roomName_out, then follow an outgoing flow edge
			curOut := fn.nodeOut(roomName)
			moved := false
			for j, ne := range fn.adj[curOut] {
				// Skip the internal reverse edge back to _in
				if ne.to == fn.nodeIn(roomName) {
					continue
				}
				if ne.flow > 0 {
					fn.adj[curOut][j].flow--
					fn.adj[ne.to][ne.rev].flow++
					cur = ne.to
					moved = true
					break
				}
			}
			if !moved {
				return nil, fmt.Errorf("ERROR: flow extraction failed — no outgoing flow from %s", roomName)
			}
		}

		paths = append(paths, path)
	}

	return paths, nil
}

// ── Path Selection ────────────────────────────────────────────────────────────

// selectBestPaths chooses the subset of paths that minimises turns for numAnts.
// Tries adding one path at a time and stops when more paths stop helping.
func selectBestPaths(paths [][]string, numAnts int) [][]string {
	bestPaths := [][]string{paths[0]}
	bestTurns := calculateTurns(bestPaths, numAnts)

	for i := 1; i < len(paths); i++ {
		candidate := paths[:i+1]
		turns := calculateTurns(candidate, numAnts)
		if turns < bestTurns {
			bestTurns = turns
			bestPaths = make([][]string, len(candidate))
			copy(bestPaths, candidate)
		} else {
			break
		}
	}

	return bestPaths
}

// calculateTurns computes the minimum turns to move numAnts through paths.
//
// We use a greedy assignment: repeatedly assign the next ant to whichever
// path currently has the smallest "finish time" = path_len + ants_assigned.
func calculateTurns(paths [][]string, numAnts int) int {
	k := len(paths)
	lengths := make([]int, k)
	for i, p := range paths {
		lengths[i] = len(p) - 1
	}

	antsOnPath := make([]int, k)
	for i := 0; i < numAnts; i++ {
		minFinish := -1
		minIdx := 0
		for j := 0; j < k; j++ {
			ft := lengths[j] + antsOnPath[j]
			if minFinish == -1 || ft < minFinish {
				minFinish = ft
				minIdx = j
			}
		}
		antsOnPath[minIdx]++
	}

	maxTurns := 0
	for i := 0; i < k; i++ {
		if antsOnPath[i] > 0 {
			ft := lengths[i] + antsOnPath[i] - 1
			if ft > maxTurns {
				maxTurns = ft
			}
		}
	}
	return maxTurns
}

// distributeAnts returns the number of ants to assign to each path.
func distributeAnts(numAnts int, paths [][]string) []int {
	k := len(paths)
	lengths := make([]int, k)
	for i, p := range paths {
		lengths[i] = len(p) - 1
	}
	antsOnPath := make([]int, k)
	for i := 0; i < numAnts; i++ {
		minFinish := -1
		minIdx := 0
		for j := 0; j < k; j++ {
			ft := lengths[j] + antsOnPath[j]
			if minFinish == -1 || ft < minFinish {
				minFinish = ft
				minIdx = j
			}
		}
		antsOnPath[minIdx]++
	}
	return antsOnPath
}

// ── Ant Simulation ────────────────────────────────────────────────────────────

// SimulateAnts produces the turn-by-turn movement lines.
//
// Each path carries ants like a train — the first ant enters on turn 1,
// the second on turn 2, etc. Each ant moves exactly one room per turn.
//
// Example: 3 ants on path [start, A, B, end]:
//
//	Turn 1: Ant1->A
//	Turn 2: Ant1->B,   Ant2->A
//	Turn 3: Ant1->end, Ant2->B,   Ant3->A
//	Turn 4:            Ant2->end, Ant3->B
//	Turn 5:                       Ant3->end
func SimulateAnts(numAnts int, paths [][]string) []string {
	antsPerPath := distributeAnts(numAnts, paths)

	// antEntry tracks one ant's assignment
	type antEntry struct {
		id     int
		path   []string
		antIdx int // 0-based position in this path's ant queue
	}

	var allAnts []antEntry
	antID := 1
	for pi, count := range antsPerPath {
		for ai := 0; ai < count; ai++ {
			allAnts = append(allAnts, antEntry{
				id:     antID,
				path:   paths[pi],
				antIdx: ai,
			})
			antID++
		}
	}

	maxTurns := calculateTurns(paths, numAnts)
	var outputLines []string

	for turn := 1; turn <= maxTurns; turn++ {
		var moves []string

		for _, ant := range allAnts {
			// Ant enters path on turn (ant.antIdx + 1).
			// On turn T, it is at position (T - ant.antIdx) in its path.
			pos := turn - ant.antIdx
			if pos <= 0 || pos >= len(ant.path) {
				continue
			}
			roomName := ant.path[pos]
			moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, roomName))
		}

		if len(moves) > 0 {
			line := moves[0]
			for _, m := range moves[1:] {
				line += " " + m
			}
			outputLines = append(outputLines, line)
		}
	}

	return outputLines
}

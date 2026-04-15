// Package colony defines the core data structures for the lem-in project.
// Think of this as the "blueprint" for everything: rooms, tunnels, and the colony itself.
package colony

// Room represents a single room in the ant colony.
// Each room has a unique name, coordinates, and knows its neighbours.
type Room struct {
	Name    string   // Unique identifier (e.g. "start", "0", "roomA")
	X, Y    int      // Coordinates (used for visualizer bonus)
	IsStart bool     // True if this is ##start
	IsEnd   bool     // True if this is ##end
	Links   []string // Names of rooms this room is connected to
}

// Colony holds all the data parsed from the input file.
// It is the central structure passed between the parser and solver.
type Colony struct {
	NumAnts   int              // How many ants we need to move
	Rooms     map[string]*Room // All rooms, keyed by name for fast lookup
	StartRoom string           // Name of the ##start room
	EndRoom   string           // Name of the ##end room
}

// NewColony creates an empty Colony ready to be populated by the parser.
func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

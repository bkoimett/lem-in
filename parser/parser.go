// Package parser handles reading and validating the lem-in input file.
//
// The input format is:
//   <number_of_ants>
//   ##start
//   <room definitions>  (name x y)
//   ##end
//   <link definitions>  (room1-room2)
//
// Lines beginning with # are comments and are ignored (except ##start and ##end).
package parser

import (
	"bufio"
	"fmt"
	"lemin/colony"
	"os"
	"strconv"
	"strings"
)

// ParseFile reads the given file, validates it, and returns:
//   - a populated Colony struct
//   - the raw file content as a string (needed for output)
//   - an error if the file is invalid
func ParseFile(filename string) (*colony.Colony, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", fmt.Errorf("ERROR: invalid data format, cannot open file")
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("ERROR: invalid data format, error reading file")
	}

	if len(lines) == 0 {
		return nil, "", fmt.Errorf("ERROR: invalid data format, empty file")
	}

	return parseLines(lines)
}

// ParseLines processes the slice of lines from the input file.
// It is exported so it can be called directly in tests.
func ParseLines(lines []string) (*colony.Colony, string, error) {
	return parseLines(lines)
}

// parseLines is the internal implementation used by both ParseFile and ParseLines.
func parseLines(lines []string) (*colony.Colony, string, error) {
	c := colony.NewColony()

	// Track which line we are currently reading
	lineIndex := 0

	// ── Step 1: Read number of ants (must be the first non-empty line) ──
	numAnts, err := strconv.Atoi(strings.TrimSpace(lines[lineIndex]))
	if err != nil || numAnts <= 0 {
		return nil, "", fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}
	c.NumAnts = numAnts
	lineIndex++

	// ── Step 2: Read rooms and links ──
	// We use a flag to know whether the NEXT room we read is ##start or ##end
	nextIsStart := false
	nextIsEnd := false

	// Track rooms we have seen to detect duplicates
	seenRooms := make(map[string]bool)

	// We will also store all links here to validate them after rooms are known
	type rawLink struct{ a, b string }
	var pendingLinks []rawLink

	// Phase tracks what we are reading: "rooms" first, then "links"
	// In practice rooms and links can be interleaved, so we handle both.
	for ; lineIndex < len(lines); lineIndex++ {
		line := strings.TrimSpace(lines[lineIndex])

		// Skip empty lines
		if line == "" {
			continue
		}

		// ── Handle special commands ──
		if line == "##start" {
			if c.StartRoom != "" {
				return nil, "", fmt.Errorf("ERROR: invalid data format, duplicate ##start")
			}
			nextIsStart = true
			continue
		}
		if line == "##end" {
			if c.EndRoom != "" {
				return nil, "", fmt.Errorf("ERROR: invalid data format, duplicate ##end")
			}
			nextIsEnd = true
			continue
		}

		// Skip regular comments (lines starting with # but not ## commands)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// ── Is this line a link? (contains '-' but no spaces) ──
		if isLink(line) {
			parts := strings.SplitN(line, "-", 2)
			pendingLinks = append(pendingLinks, rawLink{parts[0], parts[1]})
			continue
		}

		// ── Otherwise it must be a room definition: "name x y" ──
		room, err := parseRoom(line)
		if err != nil {
			return nil, "", fmt.Errorf("ERROR: invalid data format, bad room: %s", line)
		}

		// Room names cannot start with 'L' or '#'
		if strings.HasPrefix(room.Name, "L") || strings.HasPrefix(room.Name, "#") {
			return nil, "", fmt.Errorf("ERROR: invalid data format, room name cannot start with L or #")
		}

		// Detect duplicate rooms
		if seenRooms[room.Name] {
			return nil, "", fmt.Errorf("ERROR: invalid data format, duplicate room: %s", room.Name)
		}
		seenRooms[room.Name] = true

		// Apply ##start / ##end flags
		if nextIsStart {
			room.IsStart = true
			c.StartRoom = room.Name
			nextIsStart = false
		}
		if nextIsEnd {
			room.IsEnd = true
			c.EndRoom = room.Name
			nextIsEnd = false
		}

		c.Rooms[room.Name] = room
	}

	// ── Step 3: Validate start and end exist ──
	if c.StartRoom == "" {
		return nil, "", fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if c.EndRoom == "" {
		return nil, "", fmt.Errorf("ERROR: invalid data format, no end room found")
	}

	// ── Step 4: Apply links ──
	// Track duplicates: a pair (a,b) should only appear once
	seenLinks := make(map[string]bool)

	for _, lk := range pendingLinks {
		if _, ok := c.Rooms[lk.a]; !ok {
			return nil, "", fmt.Errorf("ERROR: invalid data format, link references unknown room: %s", lk.a)
		}
		if _, ok := c.Rooms[lk.b]; !ok {
			return nil, "", fmt.Errorf("ERROR: invalid data format, link references unknown room: %s", lk.b)
		}
		if lk.a == lk.b {
			return nil, "", fmt.Errorf("ERROR: invalid data format, room links to itself: %s", lk.a)
		}

		// Normalise key so a-b and b-a are the same
		key := lk.a + "-" + lk.b
		if lk.b < lk.a {
			key = lk.b + "-" + lk.a
		}
		if seenLinks[key] {
			return nil, "", fmt.Errorf("ERROR: invalid data format, duplicate link: %s-%s", lk.a, lk.b)
		}
		seenLinks[key] = true

		// Add bidirectional connection
		c.Rooms[lk.a].Links = append(c.Rooms[lk.a].Links, lk.b)
		c.Rooms[lk.b].Links = append(c.Rooms[lk.b].Links, lk.a)
	}

	// Build the raw input string (what we print before the moves)
	rawInput := buildRawOutput(lines)

	return c, rawInput, nil
}

// parseRoom parses a single "name x y" line into a Room struct.
func parseRoom(line string) (*colony.Room, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected 3 fields, got %d", len(parts))
	}

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("x coordinate not an integer")
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("y coordinate not an integer")
	}

	return &colony.Room{
		Name: parts[0],
		X:    x,
		Y:    y,
	}, nil
}

// isLink returns true if the line looks like a tunnel definition (e.g. "0-4").
// A link has exactly one '-' and no spaces.
func isLink(line string) bool {
	if strings.Contains(line, " ") {
		return false
	}
	parts := strings.SplitN(line, "-", 2)
	if len(parts) != 2 {
		return false
	}
	// Both sides must be non-empty
	return parts[0] != "" && parts[1] != ""
}

// buildRawOutput reassembles the original file lines (stripping trailing blank
// lines) into the string we print at the start of our output.
func buildRawOutput(lines []string) string {
	var sb strings.Builder
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Package parser_test contains unit tests for the parser package.
// Run with: go test ./tests/ -v
package tests

import (
	"lemin/parser"
	"strings"
	"testing"
)

// ── Helper ────────────────────────────────────────────────────────────────────

// parseString is a test helper that parses a raw multi-line string.
func parseString(input string) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	_, _, err := parser.ParseLines(lines)
	return err
}

func parseStringOK(input string) bool {
	return parseString(input) == nil
}

// ── Valid Input Tests ─────────────────────────────────────────────────────────

func TestValidSimple(t *testing.T) {
	input := `3
##start
A 0 0
##end
B 1 0
A-B`
	if !parseStringOK(input) {
		t.Error("expected valid input to parse without error")
	}
}

func TestValidWithComments(t *testing.T) {
	input := `2
#this is a comment
##start
A 0 0
#another comment
##end
B 5 0
A-B`
	if !parseStringOK(input) {
		t.Error("comments should be ignored")
	}
}

func TestValidRoomsBeforeStart(t *testing.T) {
	// Rooms can appear before ##start in the file
	input := `2
middle 5 5
##start
A 0 0
##end
B 10 0
A-middle
middle-B`
	if !parseStringOK(input) {
		t.Error("rooms before ##start should be valid")
	}
}

func TestValidNumericRoomNames(t *testing.T) {
	input := `1
##start
0 0 0
##end
1 5 0
0-1`
	if !parseStringOK(input) {
		t.Error("numeric room names should be valid")
	}
}

// ── Invalid Input Tests ───────────────────────────────────────────────────────

func TestInvalidZeroAnts(t *testing.T) {
	input := `0
##start
A 0 0
##end
B 1 0
A-B`
	if parseStringOK(input) {
		t.Error("0 ants should be invalid")
	}
}

func TestInvalidNegativeAnts(t *testing.T) {
	input := `-5
##start
A 0 0
##end
B 1 0
A-B`
	if parseStringOK(input) {
		t.Error("negative ant count should be invalid")
	}
}

func TestInvalidNoStart(t *testing.T) {
	input := `3
A 0 0
##end
B 1 0
A-B`
	if parseStringOK(input) {
		t.Error("missing ##start should be invalid")
	}
}

func TestInvalidNoEnd(t *testing.T) {
	input := `3
##start
A 0 0
B 1 0
A-B`
	if parseStringOK(input) {
		t.Error("missing ##end should be invalid")
	}
}

func TestInvalidDuplicateRoom(t *testing.T) {
	input := `3
##start
A 0 0
A 1 1
##end
B 5 0
A-B`
	if parseStringOK(input) {
		t.Error("duplicate room name should be invalid")
	}
}

func TestInvalidLinkToUnknownRoom(t *testing.T) {
	input := `3
##start
A 0 0
##end
B 1 0
A-C`
	if parseStringOK(input) {
		t.Error("link to unknown room should be invalid")
	}
}

func TestInvalidRoomNameStartsWithL(t *testing.T) {
	input := `3
##start
Lbad 0 0
##end
B 1 0
Lbad-B`
	if parseStringOK(input) {
		t.Error("room name starting with L should be invalid")
	}
}

func TestInvalidSelfLink(t *testing.T) {
	input := `3
##start
A 0 0
##end
B 1 0
A-A
A-B`
	if parseStringOK(input) {
		t.Error("self-linking room should be invalid")
	}
}

func TestInvalidDuplicateLink(t *testing.T) {
	input := `3
##start
A 0 0
##end
B 1 0
A-B
A-B`
	if parseStringOK(input) {
		t.Error("duplicate link should be invalid")
	}
}

func TestInvalidNonIntegerCoordinates(t *testing.T) {
	input := `3
##start
A 0.5 0
##end
B 1 0
A-B`
	if parseStringOK(input) {
		t.Error("non-integer coordinates should be invalid")
	}
}

func TestInvalidEmptyFile(t *testing.T) {
	err := parseString("")
	if err == nil {
		t.Error("empty file should be invalid")
	}
}

// ── Content Tests ─────────────────────────────────────────────────────────────

func TestCorrectAntCount(t *testing.T) {
	input := `7
##start
A 0 0
##end
B 1 0
A-B`
	lines := strings.Split(strings.TrimSpace(input), "\n")
	c, _, err := parser.ParseLines(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.NumAnts != 7 {
		t.Errorf("expected 7 ants, got %d", c.NumAnts)
	}
}

func TestStartAndEndRoomsSet(t *testing.T) {
	input := `1
##start
mystart 0 0
##end
myend 10 0
mystart-myend`
	lines := strings.Split(strings.TrimSpace(input), "\n")
	c, _, err := parser.ParseLines(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.StartRoom != "mystart" {
		t.Errorf("expected StartRoom=mystart, got %s", c.StartRoom)
	}
	if c.EndRoom != "myend" {
		t.Errorf("expected EndRoom=myend, got %s", c.EndRoom)
	}
}

func TestLinksAreSymmetric(t *testing.T) {
	input := `1
##start
A 0 0
B 5 0
##end
C 10 0
A-B
B-C`
	lines := strings.Split(strings.TrimSpace(input), "\n")
	c, _, err := parser.ParseLines(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A should have B in its links, and B should have A
	hasLink := func(from, to string) bool {
		for _, l := range c.Rooms[from].Links {
			if l == to {
				return true
			}
		}
		return false
	}
	if !hasLink("A", "B") || !hasLink("B", "A") {
		t.Error("links should be stored in both directions")
	}
}

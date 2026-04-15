// internal/parser/parser_test.go
package parser

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParseValidAnts(t *testing.T) {
    content := `3
##start
room1 0 0
##end
room2 5 5
room1-room2`

    tmpFile := createTempFile(t, content)
    defer os.Remove(tmpFile)

    farm, err := ParseFile(tmpFile)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if farm.Ants != 3 {
        t.Errorf("Expected 3 ants, got %d", farm.Ants)
    }
}

func TestParseInvalidAnts(t *testing.T) {
    testCases := []struct {
        name    string
        content string
    }{
        {"negative ants", "-5\n##start\nroom 0 0\n##end\nroom2 5 5"},
        {"zero ants", "0\n##start\nroom 0 0\n##end\nroom2 5 5"},
        {"non-numeric ants", "abc\n##start\nroom 0 0\n##end\nroom2 5 5"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            tmpFile := createTempFile(t, tc.content)
            defer os.Remove(tmpFile)

            _, err := ParseFile(tmpFile)
            if err == nil {
                t.Error("Expected error, got nil")
            }
        })
    }
}

func createTempFile(t *testing.T, content string) string {
    tmpFile := filepath.Join(t.TempDir(), "test.txt")
    err := os.WriteFile(tmpFile, []byte(content), 0644)
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    return tmpFile
}
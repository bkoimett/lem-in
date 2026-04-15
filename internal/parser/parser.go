// internal/parser/parser.go
package parser

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Room struct {
    Name string
    X, Y int
}

type Farm struct {
    Ants       int
    Rooms      map[string]Room
    Start, End string
    Links      map[string]map[string]bool
}

func ParseFile(filename string) (*Farm, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    farm := &Farm{
        Rooms: make(map[string]Room),
        Links: make(map[string]map[string]bool),
    }

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        line := scanner.Text()
        lineNum++

        // Skip comments (but handle ##start/##end)
        if strings.HasPrefix(line, "#") {
            continue
        }

        if lineNum == 1 {
            ants, err := strconv.Atoi(line)
            if err != nil || ants <= 0 {
                return nil, fmt.Errorf("invalid number of ants")
            }
            farm.Ants = ants
            continue
        }

        // Handle links (contains '-')
        if strings.Contains(line, "-") {
            parts := strings.Split(line, "-")
            if len(parts) != 2 {
                return nil, fmt.Errorf("invalid link format")
            }
            // Will implement link parsing later
        }
    }

    return farm, nil
}
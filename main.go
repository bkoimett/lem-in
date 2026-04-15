package main

import (
	"fmt"
	"os"

	"lemin/parser"
	"lemin/solver"
)

func main() {
	// Ensure a file argument is provided
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Step 1: Parse the input file into a colony (rooms, tunnels, ant count)
	colony, rawInput, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Step 2: Find the best set of paths using the solver
	paths, err := solver.FindPaths(colony)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Step 3: Simulate ant movement and collect move lines
	moves := solver.SimulateAnts(colony.NumAnts, paths)

	// Step 4: Print everything — raw file first, then moves
	fmt.Print(rawInput)
	fmt.Println()
	for _, line := range moves {
		fmt.Println(line)
	}
}

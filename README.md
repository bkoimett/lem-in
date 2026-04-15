# lem-in 🐜

A Go program that finds the optimal way to move **N ants** through a colony
of rooms and tunnels, from `##start` to `##end` in the fewest possible turns.

---

## Quick Start

```bash
# Build
go build -o lem-in .

# Run with an example
./lem-in examples/example00.txt

# Or with go run
go run . examples/example00.txt

# Run tests
go test ./tests/ -v
```

---

## How It Works

1. **Parse** the input file — read ant count, rooms, and tunnels
2. **Find paths** using max-flow (Edmonds-Karp on a node-split graph)
3. **Select** the optimal subset of paths for the ant count
4. **Simulate** movement — ants move in staggered "trains" along their paths
5. **Output** the original file content followed by turn-by-turn moves

### Why Max Flow?

Each room can hold only one ant at a time. Using multiple parallel,
non-overlapping paths lets more ants move simultaneously. Max flow finds
the maximum number of such vertex-disjoint paths — and the node-split trick
turns the room-capacity constraint into a standard edge-capacity problem.

See `learning/03_maxflow.md` for a full explanation.

---

## Input Format

```
<number_of_ants>
##start
<room_name> <x> <y>
[more rooms...]
##end
<room_name> <x> <y>
[links as name1-name2]
```

- Lines starting with `#` are comments (ignored), except `##start` and `##end`
- Room names cannot start with `L` or `#`
- Coordinates must be integers
- Each tunnel connects exactly two rooms

---

## Output Format

```
<original file content>

L1-roomA L2-roomB
L1-roomC L2-roomA L3-roomB
...
```

Each line is one turn. `Lx-y` means ant number x moved to room y.

---

## Examples

```
examples/
├── example00.txt    4 ants, single path         → ≤6 turns
├── example01.txt    10 ants, 3 paths             → ≤8 turns
├── example02.txt    20 ants, 2 paths             → ≤11 turns
├── example03.txt    4 ants, 2 paths              → ≤6 turns
├── example04.txt    9 ants, named rooms          → ≤6 turns
├── example05.txt    9 ants, complex network      → ≤8 turns
├── badexample00.txt no ##start                   → ERROR
├── badexample01.txt 0 ants                       → ERROR
└── badexample02.txt no path exists              → ERROR
```

---

## Project Structure

```
lem-in/
├── main.go              Entry point
├── go.mod               Go module definition
├── colony/
│   └── colony.go        Core data structures (Room, Colony)
├── parser/
│   └── parser.go        Input file reading and validation
├── solver/
│   └── solver.go        Max-flow path finding + ant simulation
├── tests/
│   ├── parser_test.go   Parser unit tests
│   └── solver_test.go   Solver unit tests
├── examples/            Sample input files
└── learning/            Step-by-step concept guides
    ├── 01_graphs.md     What is a graph?
    ├── 02_bfs.md        Breadth-First Search
    ├── 03_maxflow.md    Max flow & the node-split trick
    ├── 04_parsing.md    Parsing techniques in Go
    ├── 05_simulation.md Simulating ant movement
    └── 06_go_practices.md  Go patterns and best practices
```

---

## Error Messages

| Error | Cause |
|-------|-------|
| `ERROR: invalid data format, invalid number of ants` | Ant count ≤ 0 or not a number |
| `ERROR: invalid data format, no start room found` | Missing `##start` |
| `ERROR: invalid data format, no end room found` | Missing `##end` |
| `ERROR: invalid data format, duplicate room: X` | Room name used twice |
| `ERROR: invalid data format, duplicate link: X-Y` | Same tunnel defined twice |
| `ERROR: invalid data format, link references unknown room: X` | Link to non-existent room |
| `ERROR: invalid data format, no path between start and end` | Graph is disconnected |

---

## Learning Resources

New to the concepts? Start with the `learning/` folder:

1. **`01_graphs.md`** — nodes, edges, adjacency lists
2. **`02_bfs.md`** — how BFS finds shortest paths
3. **`03_maxflow.md`** — the key algorithm behind this project
4. **`04_parsing.md`** — reading and validating text files in Go
5. **`05_simulation.md`** — the turn-by-turn movement model
6. **`06_go_practices.md`** — Go idioms and best practices

---

## Only Standard Library

This project uses only Go's standard library (`bufio`, `fmt`, `os`,
`strconv`, `strings`) — no external packages.

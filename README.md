# lem-in - Ant Farm Simulator

## 📖 Overview

lem-in is a Go program that simulates moving ants through a colony (graph of rooms connected by tunnels) from a start room to an end room. The goal is to find the quickest path(s) and display the movement of ants turn by turn, following specific rules:

- One ant per room (except start and end rooms)
- One ant per tunnel per turn
- Ants can only move forward through tunnels
- Find the optimal combination of paths to minimize total turns

## 🎯 Learning Objectives

By building this project, you will learn:

- **Graph Algorithms**: BFS, pathfinding, multiple path optimization
- **Data Structures**: Adjacency lists, queues, graph representation
- **Parsing**: Reading and validating structured input
- **Simulation**: Turn-based movement with constraints
- **Test-Driven Development**: Writing tests before implementation
- **Error Handling**: Graceful failure with meaningful messages

## 📋 Requirements

- **Language**: Go (only standard packages allowed)
- **Input**: Text file describing ant colony
- **Output**: Original file content + ant movements in specified format
- **Team Size**: 1-4 members

## 🚀 Getting Started

### Prerequisites

```bash
# Install Go (if not already installed)
# https://golang.org/dl/

# Verify installation
go version
```

### Initial Setup

```bash
# Create project directory
mkdir lem-in
cd lem-in

# Initialize go module
go mod init lem-in

# Create basic structure
mkdir -p parser graph pathfinder simulator tests/testdata
```

## 📁 Project Structure

```
lem-in/
├── main.go                 # Entry point
├── parser/                 # Input parsing
│   ├── parser.go
│   └── parser_test.go
├── graph/                  # Graph representation
│   ├── graph.go
│   └── graph_test.go
├── pathfinder/            # Path finding algorithms
│   ├── pathfinder.go
│   ├── multi_path.go
│   └── pathfinder_test.go
├── simulator/             # Ant movement simulation
│   ├── simulator.go
│   └── simulator_test.go
├── tests/                 # Integration tests
│   ├── integration_test.go
│   └── testdata/          # Test input files
│       ├── valid/
│       └── invalid/
├── go.mod
└── README.md
```

## 🎮 Input Format

```
<number_of_ants>
##start
<room_name> <x_coord> <y_coord>
<room_name> <x_coord> <y_coord>
...
##end
<room_name> <x_coord> <y_coord>
<room_name>-<room_name>  # tunnels
<room_name>-<room_name>
```

### Example Input
```txt
3
##start
start 1 6
middle 4 8
##end
end 9 5
start-middle
middle-end
```

### Rules for Valid Input
- Room names: no spaces, don't start with `L` or `#`
- Coordinates: integers
- No duplicate rooms or tunnels
- All tunnels must connect existing rooms
- Must have exactly one `##start` and one `##end`
- Comments start with `#` (except `##start`/`##end`)

## 📤 Output Format

The program outputs:
1. The entire input file content
2. A blank line
3. Turn-by-turn ant movements

```
L<ant>-<room> L<ant>-<room> ...  # Turn 1
L<ant>-<room> L<ant>-<room> ...  # Turn 2
...
```

## 🧪 Development Approach (MVP by MVP)

Follow this step-by-step approach, building and testing each component before moving to the next:

### MVP 1: Input Parsing ✅
- [ ] Read file from command line argument
- [ ] Parse number of ants (validate integer > 0)
- [ ] Parse rooms (name, x, y) with `##start`/`##end` markers
- [ ] Parse tunnels (room1-room2)
- [ ] Basic error handling (invalid format)
- [ ] **Test**: Valid and invalid input files

### MVP 2: Graph Building ✅
- [ ] Create graph structure (adjacency list)
- [ ] Add bidirectional tunnels
- [ ] Validate no duplicate tunnels
- [ ] Validate all tunnels reference existing rooms
- [ ] **Test**: Graph construction, duplicate detection

### MVP 3: Single Path Finding ✅
- [ ] Implement BFS for shortest path
- [ ] Return path as slice of room names
- [ ] Handle case with no path
- [ ] **Test**: Various graph topologies, disconnected graphs

### MVP 4: Multiple Paths ✅
- [ ] Find all possible paths (DFS with backtracking)
- [ ] Find node-disjoint paths
- [ ] Calculate optimal path combination for given ants
- [ ] **Test**: Path combinations, disjointness verification

### MVP 5: Ant Simulation ✅
- [ ] Create Ant struct with ID, path, position
- [ ] Simulate turn-by-turn movement
- [ ] Enforce room capacity (1 ant per room)
- [ ] Enforce tunnel usage (1 ant per tunnel per turn)
- [ ] Output moves in correct format
- [ ] **Test**: Movement logic, capacity rules

### MVP 6: Complete Integration ✅
- [ ] Wire all components together
- [ ] Add comprehensive error messages
- [ ] Optimize for performance (1000+ ants)
- [ ] Final testing with provided examples
- [ ] **Test**: Full program execution

## 🧪 Testing Strategy

### Unit Tests (per package)
```bash
# Run all unit tests
go test ./... -v

# Run specific package tests
go test ./parser/... -v
go test ./graph/... -v
```

### Integration Tests
```bash
# Create test files in tests/testdata/
go test ./tests/... -v
```

### Example Test Cases
Create test files for:
- ✅ Basic linear path: `start -> A -> B -> end`
- ✅ Multiple paths: `start -> A -> end`, `start -> B -> end`
- ✅ Complex graph with cycles
- ❌ No path between start and end
- ❌ Invalid number of ants (negative, zero, text)
- ❌ Missing start or end room
- ❌ Duplicate rooms or tunnels
- ❌ Tunnel to non-existent room

## 📊 Provided Examples

The project includes example files (from the subject). Your program should produce the exact output shown:

- `example00.txt` - Basic example (6 turns for 4 ants)
- `example01.txt` - More complex (8 turns for 10 ants)
- `example02.txt` - Multiple paths (11 turns for 20 ants)
- `example03.txt` - Graph with choices (6 turns for 4 ants)
- `example04.txt` - Another topology (6 turns for 9 ants)
- `example05.txt` - Complex (8 turns for 9 ants)
- `badexample00.txt` - Invalid input
- `badexample01.txt` - Invalid input

## 🎯 Performance Requirements

- Handle 100 ants in example06: < 1.5 minutes
- Handle 1000 ants in example07: < 2.5 minutes

## 🏆 Bonus Features

- [ ] **Visualizer**: Show ants moving through colony graphically
  ```bash
  ./lem-in ant-farm.txt | ./visualizer
  ```
- [ ] **3D Visualizer**: Enhanced visualization
- [ ] **Detailed Errors**: Specific error messages
  ```
  ERROR: invalid data format, invalid number of Ants
  ERROR: invalid data format, no start room found
  ```

## 💡 Tips & Best Practices

1. **Start Simple**: Get single path working before multiple paths
2. **Test Early**: Write tests for each function as you build
3. **Validate Everything**: Don't trust input data
4. **Use Structs**: Group related data (Room, Ant, Farm)
5. **Comment Complex Logic**: Especially pathfinding algorithms
6. **Profile Performance**: Use `go test -bench` for optimization
7. **Read Examples First**: Understand expected output format

## 🐛 Common Pitfalls

- ❌ Forgetting that tunnels are bidirectional
- ❌ Allowing ants to share rooms (except start/end)
- ❌ Allowing multiple ants in same tunnel per turn
- ❌ Not handling comments correctly
- ❌ Ignoring empty lines in input
- ❌ Not validating all rooms exist before adding tunnels

## 📚 Resources

### Go Documentation
- [File I/O](https://gobyexample.com/reading-files)
- [String Manipulation](https://gobyexample.com/string-functions)
- [Structs and Methods](https://gobyexample.com/structs)

### Algorithms
- [BFS Explanation](https://en.wikipedia.org/wiki/Breadth-first_search)
- [Multiple Pathfinding](https://en.wikipedia.org/wiki/Suurballe%27s_algorithm)

### Project-Specific
- [Original lem-in (42 project)](https://github.com/01-edu/public/tree/master/subjects/lem-in)

## 🔄 Development Workflow

```bash
# 1. Create a branch for your feature
git checkout -b mvp1-parsing

# 2. Write tests first (TDD)
# Edit parser/parser_test.go

# 3. Implement the feature
# Edit parser/parser.go

# 4. Run tests
go test ./parser/... -v

# 5. Commit working code
git add .
git commit -m "MVP1: Complete input parsing"

# 6. Merge and move to next MVP
git checkout main
git merge mvp1-parsing
```

## ✅ Definition of Done

Your project is complete when:
- [ ] All MVP requirements implemented
- [ ] All tests pass (unit + integration)
- [ ] Program handles all error cases gracefully
- [ ] Output matches example outputs exactly
- [ ] Performance meets requirements
- [ ] Code follows Go best practices (gofmt, golint)
- [ ] Code is well-commented and readable
- [ ] Bonus features implemented (optional)

## 🤝 Working in a Group

- Use Git branches for parallel development
- Review each other's code before merging
- Divide MVPs among team members
- Regular standups to coordinate

---

**Good luck! Start with MVP 1 (parsing) and work your way up. Test everything as you go!** 🐜🏠
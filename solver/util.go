package solver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// JsonObject is the base struct for decoding JSON files containing one or more nonogram puzzles
type JsonObject struct {
	Puzzles []Puzzle
}

// Puzzle has a name and two 2-dimensional slices of integers representing
// the constraints of the puzzle
type Puzzle struct {
	Name string
	Rows [][]int
	Cols [][]int
}

// ReadJsonPuzzleFile reads a json file and tries to parse it, returning a JsonObject with the 
// file structure.
func ReadJsonPuzzleFile(name string) (v JsonObject) {
	f, err := os.Open(name)
	defer f.Close()

	if err != nil {
		panic(err)	
	}

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	dec.Decode(&v)

	return
}

func (obj JsonObject) ListNames() {
	fmt.Println("Found the following puzzles: ")
	for _, puzzle := range obj.Puzzles {
		fmt.Printf("\t%d x %d - %s\n", len(puzzle.Rows), len(puzzle.Cols), puzzle.Name)
	}
}
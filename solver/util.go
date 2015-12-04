package solver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// JSONObject is the base struct for decoding JSON files containing one or more nonogram puzzles
type JSONObject struct {
	Puzzles []Puzzle
}

// Puzzle has a name and two 2-dimensional slices of integers representing
// the constraints of the puzzle
type Puzzle struct {
	Name string
	Rows [][]int
	Cols [][]int
}

// ReadJSONPuzzleFile reads a json file and tries to parse it, returning a JSONObject with the
// file structure.
func ReadJSONPuzzleFile(name string) (v JSONObject) {
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

// ListNames prints to stdout a list of all the puzzles found in a JSONObject
func (obj JSONObject) ListNames() {
	fmt.Println("Found the following puzzles: ")
	for _, puzzle := range obj.Puzzles {
		fmt.Printf("\t%d x %d - %s\n", len(puzzle.Rows), len(puzzle.Cols), puzzle.Name)
	}
}

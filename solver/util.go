package solver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"errors"
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
func ReadJSONPuzzleFile(name string) (puzzles JSONObject, err error) {
	f, err := os.Open(name)
	defer f.Close()

	if err != nil {
		return 
	}

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)
	// decodes the file in the puzzles struct
	dec.Decode(&puzzles)

	return
}

// ListNames prints to stdout a list of all the puzzles found in a JSONObject
func (obj JSONObject) ListNames() {
	fmt.Println("Found the following puzzles: ")
	for _, puzzle := range obj.Puzzles {
		fmt.Printf("\t%d x %d - %s\n", len(puzzle.Rows), len(puzzle.Cols), puzzle.Name)
	}
}

// GetByName retrieves a puzzle from the json object by its name.
// it will return an error if no puzzle is found for the given name
func (obj JSONObject) GetByName(name string) (p Puzzle, err error) {
	for _, puzzle := range obj.Puzzles {
		if puzzle.Name == name {
			p = puzzle
			return
		}	
	}
	err = errors.New("No puzzle found with the given name")
	return  
}
package gongram

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// base struct for decoding JSON files containing one or more nonogram puzzles
type JsonObject struct {
	Puzzles []Puzzle
}

// a single Puzzle has a name and two 2-dimensional slices of integers representing
// the constraints of the puzzle
type Puzzle struct {
	Name string
	Rows [][]int
	Cols [][]int
}

func ReadFile(name string) (v JsonObject) {
	f, err := os.Open(name)
	defer f.Close()

	if err != nil {
		panic(err)
		return
	}

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	dec.Decode(&v)

	return
}

func ListNames(obj JsonObject) {
	fmt.Println("Found the following puzzles: ")
	for _, puzzle := range obj.Puzzles {
		fmt.Printf("\t%d x %d - %s\n", len(puzzle.Rows), len(puzzle.Cols), puzzle.Name)
	}
}

func StringifyBoard(board Board) string {
	var buffer bytes.Buffer

	for _, line := range board.cells {
		buffer.WriteString("|")

		for _, cell := range line {
			switch cell {
			case full:
				buffer.WriteString(" # ")
			case marked:
				buffer.WriteString(" - ")
			case empty:
				buffer.WriteString("   ")
			}
		}
		buffer.WriteString("|\n")
	}
	return buffer.String()
}

func writeFile(name string) {
	f, err := os.Create(name)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	var v map[string]interface{}
	v = make(map[string]interface{})

	v["message"] = "Hello, JSON"

	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)

	err = enc.Encode(&v)

	w.Flush()
}

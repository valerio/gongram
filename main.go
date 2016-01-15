package main

import (
	"flag"
	"fmt"

	"github.com/sosdoc/gongram/solver"
)

var fileName = flag.String("f", "puzzles/nonogram.json", "The name of the JSON file containing puzzle definitions.")
var puzzleName = flag.String("p", "", "Name of the puzzle to solve. It has to be contained in the loaded file.")
var listNames = flag.Bool("l", false, "Displays the names in the puzzle file without solving.")

func main() {
	flag.Parse()
	jsonObj, err := solver.ReadJSONPuzzleFile(*fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	if *listNames || *puzzleName == "" {
		jsonObj.ListNames()
		return
	}

	puzzle, err := jsonObj.GetByName(*puzzleName)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Loaded puzzle:", puzzle.Name)
	s := solver.NewTreeSolver(puzzle)
	board := s.Solve()
	fmt.Println(board)
}

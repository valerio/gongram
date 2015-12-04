package main

import (
	"fmt"
	"github.com/sosdoc/gongram/solver"
)

// TODO: the main should parse args and solve based on that

func main() {
	inputFile := solver.ReadJSONPuzzleFile("puzzles/nonogram.json")
	puz := inputFile.Puzzles[0]
	
	s := solver.NewTreeSolver(puz)	
	board := s.Solve()
	fmt.Println(board)
}

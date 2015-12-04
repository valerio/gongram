package main

import (
	"fmt"
	"github.com/sosdoc/gongram/solver"
)

// TODO: the main should parse args and solve based on that

func main() {	
	inputFile := solver.ReadJsonPuzzleFile("puzzles/nonogram.json")
	puz := inputFile.Puzzles[0]
	line := make([]solver.Cell, 10)
	line[0], line[5] = 2, 2
	fmt.Printf("line: %v\n", line)
	fmt.Printf("Row: %v\n", puz.Rows[2])
	nel, _ := solver.SolveLine(puz.Rows[2], line)
	fmt.Printf("positions: %v\n", nel) // should print [1, 6]
}

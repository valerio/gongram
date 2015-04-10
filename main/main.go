package main

import (
	"fmt"
	"github.com/sosdoc/gongram"
)

func main() {
	inputFile := gongram.ReadFile("nonogram.txt")
	//gopbn.ListNames(inputFile)
	puz := inputFile.Puzzles[0]
	line := make([]gongram.Cell, 10)
	line[0], line[5] = 2, 2
	//var b gopbn.Board = gopbn.NewBoard(puz)
	fmt.Printf("line: %v\n", line)
	fmt.Printf("Row: %v\n", puz.Rows[2])
	nel, _ := gongram.Intersect(puz.Rows[2], line)
	fmt.Printf("positions: %v\n", nel) // should print [1, 6]

}

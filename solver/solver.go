package solver

import "bytes"

// The Solver interface exposes only a single method, which returns a solved Board.
type Solver interface {
	Solve() Board
}

// The Board type represents a grid of Cells, it is used to keep the state of the puzzle
// while it's being solved.
type Board [][]Cell

// Returns a textual representation of the Board
func (board Board) String() string {
	var buffer bytes.Buffer

	for _, line := range board {
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

// Cell is an enum used for values of single cells in the Board
// A cell can be in three states:
// 	- empty, when the solver has not yet made any assumption on it
//  - marked, when the solver is certain the cell is NOT part of the picture
//  - full, when the cell is part of the picture in the nonogram
type Cell int

const (
	empty Cell = iota
	full
	marked
)

// NewBoard creates a new empty board with the specified columns and rows
func NewBoard(rows int, columns int) (b Board) {
	b = make([][]Cell, rows)

	for i := range b {
		b[i] = make([]Cell, columns)
	}
	return
}

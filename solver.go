package gongram

type Solver interface {
	Solve() Board
}

type Board [][]Cell

type Cell int

// enum used for cell values
// marked is used to signal a cell is definitely empty
// empty is for cells that haven't been evaluated yet (might be full or marked)
const (
	empty Cell = iota
	full
	marked
)

func NewBoard(rows int, columns int) (b Board) {
	b = make([][]Cell, rows)

	for i := range b {
		b[i] = make([]Cell, columns)
	}
	return
}

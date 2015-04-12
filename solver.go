package gongram

type Solver interface {
	Solve() Board
}

type Board struct {
	cells [][]Cell
}

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
	b.cells = make([][]Cell, rows)

	for i := range b.cells {
		b.cells[i] = make([]Cell, columns)
	}
	return
}

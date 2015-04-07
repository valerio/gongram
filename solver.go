package gongram

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

func NewBoard(p Puzzle) (b Board) {
	b.cells = make([][]Cell, len(p.Rows))

	for i := range b.cells {
		b.cells[i] = make([]Cell, len(p.Cols))
	}

	return
}

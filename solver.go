package gopbn

type Board struct {
	cells [][]Cell
}

type Cell int

// enum used for solving the board
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

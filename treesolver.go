package gongram

// A TreeSolver solves a nonogram by trying to solve lines iteratively until the puzzle is completed
//
// it keeps track of lines that need to be solved yet in a slice of jobs, whenever one of the lines
// is modified the other affected lines are reactivated for a new pass of the line solver, this until
// no more jobs are available.
//
// if the puzzle can't be solved by just the line solver, it then picks a cell and fills it with a
// value (either full or marked), puts the two boards in a binary tree and resumes solving with the
// line solver using a depth-first strategy
type TreeSolver struct {
	puzzle     Puzzle
	board      Board
	jobs       []TreeSolverJob
	activeJobs int
	GuessCount int
}

type LineType int

const (
	row LineType = iota
	column
)

// A job is used to represent lines of the board (either rows or columns) that have to be solved
// executing a job means to call the function SolveLine for that job
type TreeSolverJob struct {
	ltype       LineType
	index       int
	line        []Cell
	constraints []int
	score       int
}

// NewTreeSolver returns a newly created Solver for the given puzzle
func NewTreeSolver(p Puzzle) *TreeSolver {
	t := TreeSolver{p}
	t.board = NewBoard(len(p.Rows), len(p.Cols))
	t.initJobs()
	t.activeJobs = p.Rows + p.Cols
	return &t
}

func (t *TreeSolver) Solve() Board {
	return NewBoard(0, 0) //TODO: not implemented yet
}

func (t *TreeSolver) row(index int) []Cell {
	return t.board[index]
}

func (t *TreeSolver) col(index int) []Cell {
	result := make([]Cell)
	for i := 0; i < len(t.puzzle.Rows); i++ {
		append(result, t.board[i][index])
	}
	return result
}

func (t *TreeSolver) score(lt LineType, index int) int {
	return 0 //TODO: not implemented yet
}

func (t *TreeSolver) initJobs() {
	t.jobs = make([]TreeSolverJob, len(t.puzzle.Rows)+len(t.puzzle.Cols))

	for i := 0; i < len(t.puzzle.Rows); i++ {
		append(t.jobs, TreeSolverJob{row, i, t.row(i), t.puzzle.Rows[i], t.score(row, i)})
	}

	for i := 0; i < len(t.puzzle.Cols); i++ {
		append(t.jobs, TreeSolverJob{column, i, t.col(i), t.puzzle.Cols[i], t.score(column, i)})
	}
}

package gongram

import (
	"reflect"
	"sort"
)

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

// LineType identifies a line (slice of Cell) as either a row or column
type LineType int

const (
	row LineType = iota
	column
)

// TreeSolverJob represents lines of the board (either rows or columns) that
// have yet to be solved
// the score can be used as a priority for job execution
type TreeSolverJob struct {
	ltype       LineType
	index       int
	line        []Cell
	constraints []int
	score       int
}

type TreeSolverJobs []TreeSolverJob

func (slice *TreeSolverJobs) Len() int {
	len(slice)
}

func (slice *TreeSolverJobs) Less(i, j int) bool {
	return slice[i].score < slice[j].score
}

func (slice *TreeSolverJobs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// NewTreeSolver returns a newly created Solver for the given puzzle
func NewTreeSolver(p Puzzle) *TreeSolver {
	t := TreeSolver{puzzle: p}
	t.board = NewBoard(len(p.Rows), len(p.Cols))
	t.initJobs()
	t.activeJobs = len(p.Rows) + len(p.Cols)
	return &t
}

func (t *TreeSolver) Solve() Board {
	return NewBoard(0, 0) //TODO: not implemented yet
}

func (t *TreeSolver) row(index int) []Cell {
	return t.board[index]
}

func (t *TreeSolver) col(index int) []Cell {
	result := make([]Cell, len(t.puzzle.Cols))
	for i := 0; i < len(t.puzzle.Rows); i++ {
		result = append(result, t.board[i][index])
	}
	return result
}

func (t *TreeSolver) setLine(lt LineType, index int, line []Cell) {
	if lt == row {
		t.board[index] = line
	} else {
		for i := 0; i < t.puzzle.Rows; i++ {
			t.board[i][index] = line[i]
		}
	}
}

func (t *TreeSolver) score(lt LineType, index int) int {
	constraints := t.puzzle[lt][index]
	l, b, n := len(t.puzzle[lt]), 0, len(constraints)

	for _, c := range constraints {
		b += c
	}

	if b == l {
		return l
	} else {
		return b*(n+1) + n*(n-l-1)
	}
}

func (t *TreeSolver) initJobs() {
	t.jobs = make([]TreeSolverJob, len(t.puzzle.Rows)+len(t.puzzle.Cols))

	for i := 0; i < len(t.puzzle.Rows); i++ {
		t.jobs = append(t.jobs, TreeSolverJob{row, i, t.row(i), t.puzzle.Rows[i], t.score(row, i)})
	}

	for i := 0; i < len(t.puzzle.Cols); i++ {
		t.jobs = append(t.jobs, TreeSolverJob{column, i, t.col(i), t.puzzle.Cols[i], t.score(column, i)})
	}
}

func (t *TreeSolver) logicSolve() (emptyCells int, ok bool) {
	ok = true

	if len(t.jobs) == 0 {
		t.initJobs()
	}

	for len(t.jobs) > 0 {
		// pop the last job from the slice
		job, t.jobs = t.jobs[len(t.jobs)-1], t.jobs[:len(t.jobs)-1]

		if newLine, success := intersect(job.constraints, job.line); !success {
			//contradiction, stops solving
			ok = false
			return
		}

		if !reflect.DeepEqual(job.line, newLine) {
			// update the line and jobs
			t.setLine(job.ltype, job.index, newLine)

		}

	}

}

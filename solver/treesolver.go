package solver

import (
	"reflect"
	"sort"
)

// LineType identifies a line (slice of Cell) as either a row or column
type LineType int

const (
	row LineType = iota
	column
)

// treeSolverJob represents lines of the board (either rows or columns) that
// have yet to be solved.
// The score is used as priority for job execution.
type treeSolverJob struct {
	ltype       LineType
	index       int
	line        []Cell
	constraints []int
	score       int
}

type treeSolverJobs []treeSolverJob

func (slice treeSolverJobs) Len() int {
	return len(slice)
}

func (slice treeSolverJobs) Less(i, j int) bool {
	return slice[i].score < slice[j].score
}

func (slice treeSolverJobs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

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
	jobs       treeSolverJobs
	activeJobs int
	GuessCount int
}

// NewTreeSolver returns a newly created Solver for the given puzzle
func NewTreeSolver(p Puzzle) *TreeSolver {
	t := TreeSolver{puzzle: p}
	t.board = NewBoard(len(p.Rows), len(p.Cols))
	t.initJobs()
	t.activeJobs = len(p.Rows) + len(p.Cols)
	return &t
}

// Solve implements the Solver interface, returning a fully solved Board
func (t *TreeSolver) Solve() Board {
	_, _ = t.logicSolve()
	return t.board //TODO: not implemented yet
}

func (t *TreeSolver) getLine(lt LineType, index int) []Cell {
	if lt == row {
		return t.board[index]
	}
	result := make([]Cell, len(t.puzzle.Cols))
	for i := 0; i < len(t.puzzle.Rows); i++ {
		result[i] = t.board[i][index]
	}
	return result
}

func (t *TreeSolver) emptyCells() int {
	count := 0
	for _, line := range t.board {
		for _, cell := range line {
			if cell == empty {
				count++
			}
		}
	}
	return count
}

func (t *TreeSolver) setLine(lt LineType, index int, line []Cell) {
	if lt == row {
		t.board[index] = line
	} else {
		for i := 0; i < len(t.puzzle.Rows); i++ {
			t.board[i][index] = line[i]
		}
	}
}

func (t *TreeSolver) score(lt LineType, index int) int {
	var constraints []int
	var l int

	if lt == row {
		constraints = t.puzzle.Rows[index]
		l = len(t.puzzle.Rows)
	} else {
		constraints = t.puzzle.Cols[index]
		l = len(t.puzzle.Cols)
	}

	b, n := 0, len(constraints)

	for _, c := range constraints {
		b += c
	}

	if b == l {
		return l
	}
	return b*(n+1) + n*(n-l-1)
}

func (t *TreeSolver) initJobs() {
	t.jobs = make([]treeSolverJob, 0)

	for i := 0; i < len(t.puzzle.Rows); i++ {
		t.jobs = append(t.jobs, treeSolverJob{row, i, t.getLine(row, i), t.puzzle.Rows[i], t.score(row, i)})
	}

	for i := 0; i < len(t.puzzle.Cols); i++ {
		t.jobs = append(t.jobs, treeSolverJob{column, i, t.getLine(column, i), t.puzzle.Cols[i], t.score(column, i)})
	}
}

func (t *TreeSolver) updateJobs(oldJob treeSolverJob, newLine []Cell) {
	count := 0
	for i, v := range newLine {
		if v != oldJob.line[i] {
			found := false
			count++

			for _, job := range t.jobs {
				if job.ltype != oldJob.ltype && job.index == i {
					// update score for all jobs of columns if oldjob is a row, rows otherwise
					job.score = t.score(job.ltype, job.index)
					found = true
				}
			}

			if !found {
				var lt LineType
				var constraints []int
				if oldJob.ltype == row {
					lt = column
					constraints = t.puzzle.Cols[i]
				} else {
					lt = row
					constraints = t.puzzle.Rows[i]
				}

				t.jobs = append(t.jobs, treeSolverJob{lt, i, t.getLine(lt, i), constraints, t.score(lt, i)})
			}
		}
	}
}

func (t *TreeSolver) logicSolve() (emptyCells int, ok bool) {
	ok = true

	if len(t.jobs) == 0 {
		t.initJobs()
	}

	for len(t.jobs) > 0 {
		// pop the last job from the slice
		var job treeSolverJob
		job, t.jobs = t.jobs[len(t.jobs)-1], t.jobs[:len(t.jobs)-1]

		newLine, success := intersect(job.constraints, job.line)

		if !success {
			//contradiction, stops solving
			ok = false
			return
		}

		if !reflect.DeepEqual(job.line, newLine) {
			// update the line and jobs
			t.setLine(job.ltype, job.index, newLine)
			t.updateJobs(job, newLine)
			sort.Sort(t.jobs)
		}
	}
	emptyCells = t.emptyCells()
	return
}

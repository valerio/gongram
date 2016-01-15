package solver

// FastLineSolver operates by determining the starting index of each of the blocks indicated in the constraints.
// if the constraints are [2,3], the solver will try to place a block of 2 and a block of 3 separated by at least one
// empty space, if possible.
// This line solver gets two possible (partial) solutions by trying to place blocks to the leftmost and rightmost
// positions, it then intersects the two to get a single (partial) solution that is guaranteed to be correct.
//
// This solver can miss some logical clues and give incomplete solutions, but it can still solve some simple puzzles
// by just iterating on each row and column. For more complex puzzles it needs to be supplemented with ulterior logic.
type FastLineSolver struct {
	state         flState
	isLeftSolver  bool
	currentIndex  int
	blockIndex    int
	backtracking  bool
	positions     []int
	coverage      []int
	contradiction bool
	constraints   []int
	line          []Cell
}

type flState int

// The possible states for the line solver functions.
// A brief description of each state
// newblock: 	start of a block, fails if line is over or checks the rest of the line if no more blocks are present
// placeblock: 	places the block cells one by one until the end is reached or the block doesn't fit
// finalspace: 	last space of current block, leaves one space blank and goes to the next block if no obstacles are found
// checkrest: 	checks that the remainder of the line doesn't have blocks already placed, then ends
// backtrack: 	goes to the previous block and tries to move it forward (advance), fails if it's the first block
// advanceblock:tries to move forward the current block in order to cover cells that are already full
const (
	halt flState = iota
	newblock
	placeblock
	finalspace
	checkrest
	backtrack
	advanceblock
)

// SolveLine accepts a list of costraints and a line (list of Cell)
// it will try to fill as many blocks possible in the line under the given
// constraints and the already filled in cells in the line
// the result is a line with at least the same amount of cells filled in
// the given input
//
// SolveLine only fails when solving the line with the given constraints
// is impossible (a contradiction is found).
// This usually means that the solver made a wrong guess at some point
// and backtracking will be needed for solving the puzzle.
func SolveLine(constraints []int, line []Cell) (result []Cell, ok bool) {
	result, ok = intersect(constraints, line)
	return
}

func intersect(constraints []int, line []Cell) (result []Cell, ok bool) {
	result = make([]Cell, len(line))
	ok = true
	// if the line is completely empty or full the solution is trivial
	if constraints[0] == 0 {
		for i := range result {
			result[i] = marked
		}
		return
	} else if constraints[0] == len(line) {
		for i := range result {
			result[i] = full
		}
		return
	}

	changed, rb, lb := 0, 0, 0
	lgap, rgap := true, true

	leftLine := make([]Cell, len(line))
	copy(leftLine, line)

	rightLine := make([]Cell, len(line))
	copy(rightLine, line)

	leftSolver := newLeftLineSolver(constraints, leftLine)
	rightSolver := newRightLineSolver(constraints, rightLine)

	left := leftSolver.solve()
	right := rightSolver.solve()

	if left == nil || right == nil {
		ok = false
		return
	}

	for i := 0; i < len(line); i++ {
		if !lgap && left[lb]+constraints[lb] == i {
			lgap = true
			lb++
		}
		if lgap && lb < len(constraints) && left[lb] == i {
			lgap = false
		}
		if !rgap && right[rb]+1 == i {
			rgap = true
			rb++
		}
		if rgap && rb < len(constraints) && right[rb]-constraints[rb]+1 == i {
			rgap = false
		}
		if lgap == rgap && lb == rb {
			if lgap {
				result[i] = marked
			} else {
				result[i] = full
			}
			changed++
		}

	}

	return
}

func newLeftLineSolver(constraints []int, line []Cell) FastLineSolver {
	return FastLineSolver{
		currentIndex:  0,
		blockIndex:    0,
		backtracking:  false,
		state:         newblock,
		positions:     make([]int, len(constraints)),
		coverage:      make([]int, len(constraints)),
		constraints:   constraints,
		line:          line,
		isLeftSolver:  true,
		contradiction: false,
	}
}

func newRightLineSolver(constraints []int, line []Cell) FastLineSolver {
	return FastLineSolver{
		currentIndex:  0,
		blockIndex:    len(constraints) - 1,
		backtracking:  false,
		state:         newblock,
		positions:     make([]int, len(constraints)),
		coverage:      make([]int, len(constraints)),
		constraints:   constraints,
		line:          line,
		isLeftSolver:  false,
		contradiction: false,
	}
}

// The solve functions implements a basic finite state machine, it executes until the solver reaches the halt state.
// If a contradiction has been found, no solution can be found and the solver stops. This means the puzzle currently
// contains an error.
func (ls *FastLineSolver) solve() []int {
	for ls.state != halt {
		switch ls.state {
		case newblock:
			ls.newBlock()
		case placeblock:
			ls.placeBlock()
		case finalspace:
			ls.finalSpace()
		case checkrest:
			ls.checkRest()
		case backtrack:
			ls.backtrack()
		case advanceblock:
			ls.advanceBlock()
		}
	}

	if ls.contradiction {
		return nil
	}
	return ls.positions
}

func (ls *FastLineSolver) newBlock() {
	if ls.blockIndexOutOfBounds() {
		if ls.isFirstBlockIndex() {
			ls.currentIndex = ls.startOfLine()
		}
		ls.blockIndex = ls.previousBlock()
		ls.state = checkrest
		return
	}

	if ls.isFirstBlockIndex() {
		ls.positions[ls.blockIndex] = ls.startOfLine()
	} else {
		ls.positions[ls.blockIndex] = ls.nextIndex()
	}

	if ls.blockPositionOutOfBounds() {
		ls.state = halt
		ls.contradiction = true
		return
	}
	ls.state = placeblock
}

func (ls *FastLineSolver) placeBlock() {
	//move block forward if line is marked
	for ls.line[ls.positions[ls.blockIndex]] == marked {
		ls.positions[ls.blockIndex] = ls.nextBlockPosition()
		if ls.blockPositionOverflow() {
			ls.state = halt
			ls.contradiction = true
			return
		}
	}

	ls.currentIndex = ls.positions[ls.blockIndex]

	if ls.line[ls.currentIndex] != full {
		ls.coverage[ls.blockIndex] = -1
	} else {
		ls.coverage[ls.blockIndex] = ls.currentIndex //we need to cover position i (cannot move ls.blockIndex past i)
	}
	ls.currentIndex = ls.nextIndex()

	for ls.blockIsSmallerThanConstraint() {

		if ls.currentIndexOutOfBounds() {
			ls.state = halt
			ls.contradiction = true
			return
		}

		if ls.line[ls.currentIndex] == marked {
			if ls.coverage[ls.blockIndex] == -1 {
				ls.positions[ls.blockIndex] = ls.currentIndex
				ls.state = placeblock
			} else {
				ls.state = backtrack
			}
			return
		}

		if ls.coverage[ls.blockIndex] == -1 && ls.line[ls.currentIndex] == 1 {
			ls.coverage[ls.blockIndex] = ls.currentIndex
		}

		ls.currentIndex = ls.nextIndex()
	}
	ls.state = finalspace
}

func (ls *FastLineSolver) finalSpace() {
	for !ls.currentIndexOutOfBounds() && ls.line[ls.currentIndex] == full {
		if ls.coverage[ls.blockIndex] == ls.positions[ls.blockIndex] {
			ls.state = backtrack
			return
		}

		ls.positions[ls.blockIndex] = ls.nextBlockPosition()

		if ls.coverage[ls.blockIndex] == -1 && ls.line[ls.currentIndex] == full {
			ls.coverage[ls.blockIndex] = ls.currentIndex
		}

		ls.currentIndex = ls.nextIndex()
	}

	if ls.backtracking && ls.coverage[ls.blockIndex] == -1 {
		ls.backtracking = false
		ls.state = advanceblock
		return
	}

	if ls.currentIndexOutOfBounds() && ls.blockIndexInBounds() {
		ls.state = halt
		ls.contradiction = true
		return
	}

	ls.blockIndex = ls.nextBlock()
	ls.backtracking = false
	ls.state = newblock
}

func (ls *FastLineSolver) checkRest() {
	//checks that remaining cells are empty or marked
	for !ls.currentIndexOutOfBounds() {
		if ls.line[ls.currentIndex] == full {
			//move the last block forward
			ls.moveBlockForward()
			ls.state = advanceblock
			return
		}
		ls.currentIndex = ls.nextIndex()
	}
	ls.state = halt
}

func (ls *FastLineSolver) backtrack() {
	ls.blockIndex = ls.previousBlock()
	if ls.blockIndexBeforeBounds() {
		ls.state = halt
		ls.contradiction = true
		return
	}
	ls.moveBlockForward()
	ls.state = advanceblock
}

func (ls *FastLineSolver) advanceBlock() {
	for ls.coverage[ls.blockIndex] < 0 && ls.isPositionOfBlockNotCovered() {
		if ls.line[ls.currentIndex] == marked {
			if ls.coverage[ls.blockIndex] > 0 {
				ls.state = backtrack
			} else {
				ls.positions[ls.blockIndex] = ls.nextIndex()
				ls.backtracking = true
				ls.state = placeblock
			}
			return
		}

		ls.positions[ls.blockIndex] = ls.nextBlockPosition()

		if ls.line[ls.currentIndex] == full {
			ls.currentIndex = ls.nextIndex()
			if ls.coverage[ls.blockIndex] == -1 {
				ls.coverage[ls.blockIndex] = ls.currentIndex - 1
			}
			ls.state = finalspace
			return
		}

		ls.currentIndex = ls.nextIndex()

		if ls.currentIndexOutOfBounds() {
			ls.state = halt
			ls.contradiction = true
			return
		}
	}
	ls.state = backtrack
}

// The following functions are used to abstract the left and right Solver.
// They encapsulate common conditions and operations that are logically the same.

func (ls *FastLineSolver) blockIndexOutOfBounds() bool {
	if ls.isLeftSolver {
		return ls.blockIndex >= len(ls.constraints)
	}
	return ls.blockIndex < 0
}

func (ls *FastLineSolver) blockPositionOutOfBounds() bool {
	if ls.isLeftSolver {
		return ls.positions[ls.blockIndex] == len(ls.line)
	}
	return ls.positions[ls.blockIndex]-ls.constraints[ls.blockIndex]+1 < 0
}

func (ls *FastLineSolver) blockPositionOverflow() bool {
	if ls.isLeftSolver {
		return ls.positions[ls.blockIndex] >= len(ls.line)
	}
	return ls.positions[ls.blockIndex] < 0
}

func (ls *FastLineSolver) isFirstBlockIndex() bool {
	if ls.isLeftSolver {
		return ls.blockIndex == 0
	}
	return ls.blockIndex == len(ls.constraints)-1
}

func (ls *FastLineSolver) startOfLine() int {
	if ls.isLeftSolver {
		return 0
	}
	return len(ls.line) - 1
}

func (ls *FastLineSolver) nextIndex() int {
	if ls.isLeftSolver {
		return ls.currentIndex + 1
	}
	return ls.currentIndex - 1
}

func (ls *FastLineSolver) previousIndex() int {
	if ls.isLeftSolver {
		return ls.currentIndex - 1
	}
	return ls.currentIndex + 1
}

func (ls *FastLineSolver) nextBlock() int {
	if ls.isLeftSolver {
		return ls.blockIndex + 1
	}
	return ls.blockIndex - 1
}

func (ls *FastLineSolver) previousBlock() int {
	if ls.isLeftSolver {
		return ls.blockIndex - 1
	}
	return ls.blockIndex + 1
}

func (ls *FastLineSolver) nextBlockPosition() int {
	if ls.isLeftSolver {
		return ls.positions[ls.blockIndex] + 1
	}
	return ls.positions[ls.blockIndex] - 1
}

func (ls *FastLineSolver) previousBlockPosition() int {
	if ls.isLeftSolver {
		return ls.positions[ls.blockIndex] - 1
	}
	return ls.positions[ls.blockIndex] + 1
}

func (ls *FastLineSolver) blockIsSmallerThanConstraint() bool {
	if ls.isLeftSolver {
		return ls.currentIndex-ls.positions[ls.blockIndex] < ls.constraints[ls.blockIndex]
	}
	return ls.positions[ls.blockIndex]-ls.currentIndex < ls.constraints[ls.blockIndex]
}

func (ls *FastLineSolver) currentIndexOutOfBounds() bool {
	if ls.isLeftSolver {
		return ls.currentIndex >= len(ls.line)
	}
	return ls.currentIndex < 0
}

func (ls *FastLineSolver) blockIndexInBounds() bool {
	if ls.isLeftSolver {
		return ls.blockIndex < len(ls.constraints)-1
	}
	return ls.blockIndex > 0
}

func (ls *FastLineSolver) moveBlockForward() {
	if ls.isLeftSolver {
		ls.currentIndex = ls.positions[ls.blockIndex] + ls.constraints[ls.blockIndex]
	} else {
		ls.currentIndex = ls.positions[ls.blockIndex] + ls.constraints[ls.blockIndex]
	}
}

func (ls *FastLineSolver) blockIndexBeforeBounds() bool {
	if ls.isLeftSolver {
		return ls.blockIndex < 0
	}
	return ls.blockIndex > len(ls.constraints)-1
}

func (ls *FastLineSolver) isPositionOfBlockNotCovered() bool {
	if ls.isLeftSolver {
		return ls.positions[ls.blockIndex] < ls.coverage[ls.blockIndex]
	}
	return ls.positions[ls.blockIndex] > ls.coverage[ls.blockIndex]
}

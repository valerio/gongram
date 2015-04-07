package gongram

// The possible states for the line solver functions.
// the line solver operates by determining the starting index of each of the blocks indicated in the constraints
// if the constraints are [2,3], the solver will try to place a block of 2 and a block of 3 separated by at least one
// empty space, if possible.
// This line solver gets two possible (partial) solutions by trying to place blocks to the leftmost and rightmost
// positions, it then intersects the two to get a single (partial) solution that is guaranteed to be correct.
//
// This solver can miss some logical clues and give incomplete solutions, but it can still solve some simple puzzles
// by just iterating on each row and column. For more complex puzzles it needs to be supplemented with ulterior logic.
//
// A brief description of each state
// newblock: 	start of a block, fails if line is over or checks the rest of the line if no more blocks are present
// placeblock: 	places the block cells one by one until the end is reached or the block doesn't fit
// finalspace: 	last space of current block, leaves one space blank and goes to the next block if no obstacles are found
// checkrest: 	checks that the remainder of the line doesn't have blocks already placed, then ends
// backtrack: 	goes to the previous block and tries to move it forward (advance), fails if it's the first block
// advanceblock:tries to move forward the current block in order to cover cells that are already full
//
const (
	halt = iota
	newblock
	placeblock
	finalspace
	checkrest
	backtrack
	advanceblock
)

func Intersect(constraints []int, line []Cell) []Cell {
	return []Cell{empty} //TODO: not implemented yet
}

func LeftSolve(constraints []int, line []Cell) []int {
	i, block := 0, 0
	backtracking := false
	state := newblock
	positions := make([]int, len(constraints))
	coverage := make([]int, len(constraints))

Loop:
	for state != halt {
		switch state {
		case newblock:
			//set the first index of a new block
			if block >= len(constraints) {
				if block == 0 {
					i = 0
				}
				block--
				state = checkrest
				continue Loop
			}

			if block == 0 {
				positions[block] = 0
			} else {
				positions[block] = i + 1
			}

			if positions[block] == len(line) {
				return nil
			}
			state = placeblock

		case placeblock:

			//move block forward if line is marked
			for line[positions[block]] == marked {
				positions[block]++
				if positions[block] >= len(line) {
					return nil
				}
			}

			i = positions[block]

			if line[i] != full {
				coverage[block] = -1
			} else {
				coverage[block] = i //we need to cover position i (cannot move block past i)
			}
			i++

			for i-positions[block] < constraints[block] {

				if i >= len(line) {
					return nil
				}

				if line[i] == marked {
					if coverage[block] == -1 {
						positions[block] = i
						state = placeblock
					} else {
						state = backtrack
					}
					continue Loop
				}

				if coverage[block] == -1 && line[i] == 1 {
					coverage[block] = i
				}

				i++
			}
			state = finalspace

		case finalspace:

			for i < len(line) && line[i] == full {

				if coverage[block] == positions[block] {
					state = backtrack
					continue Loop
				}

				positions[block]++

				if coverage[block] == -1 && line[i] == full {
					coverage[block] = i
				}

				i++
			}

			if backtracking && coverage[block] == -1 {
				backtracking = false
				state = advanceblock
				continue Loop
			}

			if i >= len(line) && block < len(constraints)-1 {
				return nil
			}

			block++
			backtracking = false
			state = newblock

		case checkrest:
			//checks that remaining cells are empty or marked
			for i < len(line) {
				if line[i] == full {
					//move the last block forward
					i = positions[block] + constraints[block]
					state = advanceblock
					continue Loop
				}
				i++
			}
			state = halt

		case backtrack:
			block--
			if block < 0 {
				return nil
			}
			i = positions[block] + constraints[block]
			state = advanceblock

		case advanceblock:

			for coverage[block] < 0 && positions[block] < coverage[block] {
				if line[i] == marked {
					if coverage[block] > 0 {
						state = backtrack
					} else {
						positions[block] = i + 1
						backtracking = true
						state = placeblock
					}
					continue Loop
				}

				positions[block]++

				if line[i] == full {
					i++
					if coverage[block] == -1 {
						coverage[block] = i - 1
					}
					state = finalspace
					continue Loop
				}

				i++

				if i >= len(line) {
					return nil
				}
			}
			state = backtrack
		}

	}

	return positions
}

func RightSolve(constraints []int, line []Cell) []int {
	return []int{} //TODO: not implemented yet
}

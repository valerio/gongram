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

func Intersect(constraints []int, line []Cell) (result []Cell, ok bool) {
	result = make([]Cell, len(line))
	ok = true

	// if the line is completely empty or full the solution is trivial
	if constraints[0] == 0 {
		for i, _ := range result {
			result[i] = marked
		}
		return
	} else if constraints[0] == len(line) {
		for i, _ := range result {
			result[i] = full
		}
		return
	}

	changed, rb, lb := 0, 0, 0
	lgap, rgap := true, true

	left := LeftSolve(constraints, line)
	right := RightSolve(constraints, line)

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

var OutL, OutR chan []int = make(chan []int), make(chan []int)

// IntersectP acts the same as Intersect but executes the leftmost and rightmost solver functions in parallel
// then it waits for both results to be sent on their respective channels and proceeds to compute the intersection
//
// this shouldn't be faster than the regular Intersect, as the left/right solvers are pretty fast anyway and
// the overhead from communication/sleeping will probably be higher than the time saved
func IntersectP(constraints []int, line []Cell) (result []Cell, ok bool) {
	result = make([]Cell, len(line))
	ok = true

	// if the line is completely empty or full the solution is trivial
	if constraints[0] == 0 {
		for i, _ := range result {
			result[i] = marked
		}
		return
	} else if constraints[0] == len(line) {
		for i, _ := range result {
			result[i] = full
		}
		return
	}

	changed, rb, lb := 0, 0, 0
	lgap, rgap := true, true

	go func() {
		OutL <- LeftSolve(constraints, line)
	}()

	go func() {
		OutR <- RightSolve(constraints, line)
	}()

	left := <-OutL
	right := <-OutR

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
	i, block, maxblock := 0, len(constraints)-1, len(constraints)-1
	backtracking := false
	state := newblock
	positions := make([]int, len(constraints))
	coverage := make([]int, len(constraints))
Loop:
	for state != halt {
		switch state {
		case newblock:
			if block < 0 {
				if block == maxblock {
					i = len(line) - 1
				}
				block++
				state = checkrest
				continue Loop
			}

			if block == maxblock {
				positions[block] = len(line) - 1
			} else {
				positions[block] = i - 1
			}

			if positions[block]-constraints[block]+1 < 0 {
				return nil
			}
			state = placeblock

		case placeblock:
			for line[positions[block]] == marked {
				positions[block]--
				if positions[block] < 0 {
					return nil
				}
			}

			i = positions[block]

			if line[i] != full {
				coverage[block] = -1
			} else {
				coverage[block] = i
			}
			i--

			for positions[block]-i < constraints[block] {

				if i < 0 {
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

				i--
			}
			state = finalspace

		case finalspace:
			for i >= 0 && line[i] == full {

				if coverage[block] == positions[block] {
					state = backtrack
					continue Loop
				}

				positions[block]--

				if coverage[block] == -1 && line[i] == full {
					coverage[block] = i
				}

				i--
			}

			if backtracking && coverage[block] == -1 {
				backtracking = false
				state = advanceblock
				continue Loop
			}

			if i < 0 && block > 0 {
				return nil
			}

			block--
			backtracking = false
			state = newblock

		case checkrest:
			for i >= 0 {
				if line[i] == full {
					i = positions[block] - constraints[block]
					state = advanceblock
					continue Loop
				}
				i--
			}
			state = halt

		case backtrack:
			block++
			if block > maxblock {
				return nil
			}
			i = positions[block] - constraints[block]
			state = advanceblock

		case advanceblock:
			for coverage[block] < 0 && positions[block] > coverage[block] {
				if line[i] == marked {
					if coverage[block] > 0 {
						state = backtrack
					} else {
						positions[block] = i - 1
						backtracking = true
						state = placeblock
					}
					continue Loop
				}

				positions[block]--

				if line[i] == full {
					i--
					if coverage[block] == -1 {
						coverage[block] = i - 1
					}
					state = finalspace
					continue Loop
				}

				i--

				if i < 0 {
					return nil
				}
			}
			state = backtrack
		}
	}
	return positions
}

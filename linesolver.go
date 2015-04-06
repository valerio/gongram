package gopbn

import "fmt"

const (
	halt = iota
	newblock
	placeblock
	finalspace
	checkrest
	backtrack
	advanceblock
)

func intersect() {

}

func log(in interface{}) {
	fmt.Printf("\tlog:%v\n", in)
}

func LeftSolve(constraints []int, line []Cell) []int {
	i, block := 0, 0
	backtracking := false
	state := newblock
	positions := make([]int, len(constraints))
	coverage := make([]int, len(constraints))

Loop:
	for state != halt {
		//log(state)
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
				coverage[block] = i
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

func rightSolve() {

}

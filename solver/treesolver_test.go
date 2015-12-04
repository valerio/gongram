package solver

import "testing"

func TestSolve(t *testing.T) {
	inputFile := ReadJSONPuzzleFile("../puzzles/nonogram.json")
	puz := inputFile.Puzzles[0]
	if puz.Name != "mushroom" {

	}
	s := NewTreeSolver(puz)
	board := s.Solve()

	if board[9][9] != marked {
		t.FailNow()
	}
}

func TestReadPuzzleFile(t *testing.T) {
	inputFile := ReadJSONPuzzleFile("../puzzles/nonogram.json")
	puz := inputFile.Puzzles[0]
	if puz.Name != "mushroom" {
		t.FailNow()
	}

	if len(puz.Rows) != 10 {
		t.FailNow()
	}
}

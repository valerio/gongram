package solver

import (
	"reflect"
	"testing"
)

func TestLeftSolve(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 2, 0, 0, 0}
	constraints := []int{3, 3}
	expected := []int{0, 5}

	solver := newLeftLineSolver(constraints, line)
	result := solver.solve()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 1, 0, 1, 0, 2, 0, 0}
	constraints = []int{1, 2}
	expected = []int{1, 3}

	solver = newLeftLineSolver(constraints, line)
	result = solver.solve()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}

func TestRightSolve(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 2, 0, 0, 0}
	constraints := []int{3, 3}
	expected := []int{3, 7}

	solver := newRightLineSolver(constraints, line)
	result := solver.solve()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 0, 0, 0, 0, 2, 0, 0}
	constraints = []int{1, 3}
	expected = []int{0, 4}

	solver = newRightLineSolver(constraints, line)
	result = solver.solve()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}

func TestIntersect(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 0}
	constraints := []int{5}
	expected := []Cell{1, 1, 1, 1, 1}

	result, _ := intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 0, 0, 0, 0}
	constraints = []int{0}
	expected = []Cell{2, 2, 2, 2, 2}

	result, _ = intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}	
		
	line = []Cell{0, 0, 0, 0, 0}
	constraints = []int{3}
	expected = []Cell{0, 0, 1, 0, 0}

	result, _ = intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
	
	line = []Cell{1, 0, 1, 0, 0}
	constraints = []int{3}
	expected = []Cell{1, 1, 1, 2, 2}

	result, _ = intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}

func BenchmarkIntersect(b *testing.B) {
	line := []Cell{0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 1, 0}
	constraints := []int{3, 3, 1, 4, 2}

	for i := 0; i < b.N; i++ {
		intersect(constraints, line)
	}
}

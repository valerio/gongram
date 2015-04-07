package gongram

import (
	"reflect"
	"testing"
)

func TestLeftSolve(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 2, 0, 0, 0}
	constraints := []int{3, 3}
	expected := []int{0, 5}

	result := LeftSolve(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 1, 0, 1, 0, 2, 0, 0}
	constraints = []int{1, 2}
	expected = []int{1, 3}

	result = LeftSolve(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}

func TestRightSolve(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 2, 0, 0, 0}
	constraints := []int{3, 3}
	expected := []int{3, 7}

	result := RightSolve(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 0, 0, 0, 0, 2, 0, 0}
	constraints = []int{1, 3}
	expected = []int{0, 4}

	result = LeftSolve(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}

func TestIntersect(t *testing.T) {
	line := []Cell{0, 0, 0, 0, 0}
	constraints := []int{5}
	expected := []Cell{1, 1, 1, 1, 1}

	result := Intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}

	line = []Cell{0, 0, 0, 0, 0}
	constraints = []int{0}
	expected = []Cell{2, 2, 2, 2, 2}

	result = Intersect(constraints, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
		t.FailNow()
	}
}
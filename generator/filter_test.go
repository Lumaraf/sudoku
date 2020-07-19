package generator

import (
	"fmt"
	"github.com/lumaraf/sudoku-checker/grid"
	"testing"
)

func TestFilter_UniqueGroup(t *testing.T) {
	f := NewFilter()
	f.Restrict(grid.GetCoordinate(0, 0), NewValueMask(1, 2))
	f.Restrict(grid.GetCoordinate(1, 0), NewValueMask(2, 3))
	f.Restrict(grid.GetCoordinate(2, 0), NewValueMask(1, 3))
	f.Restrict(grid.GetCoordinate(8, 0), NewValueMask(1, 2, 3, 4))

	solvable := f.UniqueGroup(
		grid.GetCoordinate(0, 0),
		grid.GetCoordinate(1, 0),
		grid.GetCoordinate(2, 0),
		grid.GetCoordinate(3, 0),
		grid.GetCoordinate(4, 0),
		grid.GetCoordinate(5, 0),
		grid.GetCoordinate(6, 0),
		grid.GetCoordinate(7, 0),
		grid.GetCoordinate(8, 0),
	)

	fmt.Println(f, solvable)
}

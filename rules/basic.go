package rules

import (
	"github.com/lumaraf/sudoku-checker/generator"
	"github.com/lumaraf/sudoku-checker/grid"
)

type RowRule struct{}

func (r RowRule) Filter(filter *generator.Filter) bool {
	for row := 0; row < 9; row++ {
		if !filter.UniqueGroup(
			grid.GetCoordinate(row, 0),
			grid.GetCoordinate(row, 1),
			grid.GetCoordinate(row, 2),
			grid.GetCoordinate(row, 3),
			grid.GetCoordinate(row, 4),
			grid.GetCoordinate(row, 5),
			grid.GetCoordinate(row, 6),
			grid.GetCoordinate(row, 7),
			grid.GetCoordinate(row, 8),
		) {
			return false
		}
	}
	return true
}

func (r RowRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for col := 0; col < 9; col++ {
		c := grid.GetCoordinate(current.Row(), col)
		if !state.Block(c, value) {
			return
		}
	}
	next(state)
}

type ColumnRule struct{}

func (r ColumnRule) Filter(filter *generator.Filter) bool {
	for col := 0; col < 9; col++ {
		if !filter.UniqueGroup(
			grid.GetCoordinate(0, col),
			grid.GetCoordinate(1, col),
			grid.GetCoordinate(2, col),
			grid.GetCoordinate(3, col),
			grid.GetCoordinate(4, col),
			grid.GetCoordinate(5, col),
			grid.GetCoordinate(6, col),
			grid.GetCoordinate(7, col),
			grid.GetCoordinate(8, col),
		) {
			return false
		}
	}
	return true
}

func (r ColumnRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for row := 0; row < 9; row++ {
		c := grid.GetCoordinate(row, current.Col())
		if !state.Block(c, value) {
			return
		}
	}
	next(state)
}

type SquareRule struct{}

func (r SquareRule) Filter(filter *generator.Filter) bool {
	for row := 0; row < 9; row += 3 {
		for col := 0; col < 9; col += 3 {
			if !filter.UniqueGroup(
				grid.GetCoordinate(row, col),
				grid.GetCoordinate(row, col+1),
				grid.GetCoordinate(row, col+2),
				grid.GetCoordinate(row+1, col),
				grid.GetCoordinate(row+1, col+1),
				grid.GetCoordinate(row+1, col+2),
				grid.GetCoordinate(row+2, col),
				grid.GetCoordinate(row+2, col+1),
				grid.GetCoordinate(row+2, col+2),
			) {
				return false
			}
		}
	}
	return true
}

func (r SquareRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	row := (current.Row() / 3) * 3
	col := (current.Col() / 3) * 3
	for rowOffset := 0; rowOffset < 3; rowOffset++ {
		for colOffset := 0; colOffset < 3; colOffset++ {
			if !state.Block(grid.GetCoordinate(row+rowOffset, col+colOffset), value) {
				return
			}
		}
	}
	next(state)
}

type GivenValuesRule map[grid.Coordinate]uint8

func (r GivenValuesRule) Filter(filter *generator.Filter) bool {
	for coordinate, value := range r {
		if !filter.Restrict(coordinate, generator.NewValueMask(value)) {
			return false
		}
	}
	return true
}

func (r GivenValuesRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	// nothing to do here, this rule only needs Init
	next(state)
}

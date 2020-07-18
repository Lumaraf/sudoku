package rules

import (
	"github.com/lumaraf/sudoku-checker/generator"
	"github.com/lumaraf/sudoku-checker/grid"
)

type RowRule struct{}

func (r RowRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for row := 0; row < 9; row++ {
		c := grid.GetCoordinate(row, current.Col())
		if !state.Block(c, value) {
			return
		}
	}
	next(state)
}

type ColumnRule struct{}

func (r ColumnRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for col := 0; col < 9; col++ {
		c := grid.GetCoordinate(current.Row(), col)
		if !state.Block(c, value) {
			return
		}
	}
	next(state)
}

type SquareRule struct{}

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

func (r GivenValuesRule) PreMask(maskGrid *generator.ValueMaskGrid) {
	for coordinate, value := range r {
		maskGrid.Restrict(coordinate, generator.NewValueMask(value))
	}
}

func (r GivenValuesRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	// nothing to do here, this rule only needs Init
	next(state)
}

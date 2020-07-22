package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
)

type NonConsecutiveRule struct{}

func (r NonConsecutiveRule) Set(current grid.Coordinate, value uint8, state generator.GeneratorState, next generator.NextFunc) {
	mask := generator.NewValueMask(value-1, value+1)
	if current.Row() > 0 && !state.BlockAll(grid.GetCoordinate(current.Row()-1, current.Col()), mask) {
		return
	}
	if current.Row() < 8 && !state.BlockAll(grid.GetCoordinate(current.Row()+1, current.Col()), mask) {
		return
	}
	if current.Col() > 0 && !state.BlockAll(grid.GetCoordinate(current.Row(), current.Col()-1), mask) {
		return
	}
	if current.Col() < 8 && !state.BlockAll(grid.GetCoordinate(current.Row(), current.Col()+1), mask) {
		return
	}
	next(state)
}

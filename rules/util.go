package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
)

type PlacementSet [9]int

func (p PlacementSet) Place(initialState generator.GeneratorState, moveSet [9]int, value uint8, final generator.NextFunc) {
	// fill cells greedily before stepping into final
	// waiting for everything to fall into valid places takes far too long otherwise
	row := 0
	var next generator.NextFunc
	next = func(state *generator.GeneratorState) {
		row++
		if row > 8 {
			final(state)
			return
		}
		col := moveSet[row]
		state.WithCell(grid.GetCoordinate(row, col), value, next)
	}
	next(&initialState)
}

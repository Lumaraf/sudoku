package rules

import (
	"github.com/lumaraf/sudoku-checker/generator"
	"github.com/lumaraf/sudoku-checker/grid"
)

type MirrorEntry struct {
	Active bool
	grid.Coordinate
}
type MirrorRule struct {
	Mirror [9][9]MirrorEntry
}

func (r *MirrorRule) Add(row1, col1, row2, col2 int) {
	if row1 < row2 || (row1 == row2 && col1 < col2) {
		r.Mirror[row1][col1] = MirrorEntry{
			Active:     true,
			Coordinate: grid.GetCoordinate(row2, col2),
		}
	} else {
		r.Mirror[row2][col2] = MirrorEntry{
			Active:     true,
			Coordinate: grid.GetCoordinate(row1, col1),
		}
	}
}

func (r MirrorRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	mirror := r.Mirror[current.Row()][current.Col()]
	if mirror.Active {
		mirrorValue := state.Get(mirror.Coordinate)
		if mirrorValue == 0 {
			state.WithCell(mirror.Coordinate, value, next)
		} else if mirrorValue == value {
			next(state)
		}
	}
	next(state)
}

type Area []grid.Coordinate

type UniqueAreaRule struct {
	areas [9][9][]Area
}

func NewUniqueAreaRule(areas []Area) *UniqueAreaRule {
	r := UniqueAreaRule{}
	for _, area := range areas {
		for _, coordinate := range area {
			buf := &r.areas[coordinate.Row()][coordinate.Col()]
			*buf = append(*buf, area)
		}
	}
	return &r
}

func (r UniqueAreaRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for _, area := range r.areas[current.Row()][current.Col()] {
		for _, other := range area {
			if !state.Block(other, value) {
				return
			}
		}
	}
	next(state)
}

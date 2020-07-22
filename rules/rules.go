package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
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

func (r MirrorRule) Set(current grid.Coordinate, value uint8, state generator.GeneratorState, next generator.NextFunc) {
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
type Areas []Area

func (a Areas) Filter(filter *generator.Filter) bool {
	for _, area := range a {
		if !filter.UniqueGroup(area...) {
			return false
		}
	}
	return true
}

func (a Areas) Set(value uint8, state generator.GeneratorState) bool {
	for _, area := range a {
		for _, other := range area {
			if !state.Block(other, value) {
				return false
			}
		}
	}
	return true
}

type UniqueAreaRule struct {
	areas     Areas
	cellAreas [9][9]Areas
}

func NewUniqueAreaRule(areas Areas) *UniqueAreaRule {
	r := UniqueAreaRule{
		areas: areas,
	}
	for _, area := range areas {
		for _, coordinate := range area {
			buf := &r.cellAreas[coordinate.Row()][coordinate.Col()]
			*buf = append(*buf, area)
		}
	}
	return &r
}

func (r UniqueAreaRule) Filter(filter *generator.Filter) bool {
	return r.areas.Filter(filter)
}

func (r UniqueAreaRule) Set(current grid.Coordinate, value uint8, state generator.GeneratorState, next generator.NextFunc) {
	if !r.cellAreas[current.Row()][current.Col()].Set(value, state) {
		return
	}
	next(state)
}

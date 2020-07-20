package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
)

type ThermometerRule struct {
	areas     Areas
	cellAreas [9][9]Areas
}

func NewThermometerRule(areas Areas) ThermometerRule {
	r := ThermometerRule{
		areas: areas,
	}
	for _, area := range areas {
		for _, coordinate := range area {
			buf := &r.cellAreas[coordinate.Row()][coordinate.Col()]
			*buf = append(*buf, area)
		}
	}
	return r
}

func (r ThermometerRule) Filter(filter *generator.Filter) bool {
	for _, area := range r.areas {
		cellRange := make([][2]uint8, len(area))

		min := uint8(1)
		for n := 0; n < len(area); n++ {
			mask := filter.Get(area[n])
			cellRange[n][0] = min
			if maskMin := mask.Minimum(); maskMin > min {
				min = maskMin
			}
			min++
		}

		max := uint8(9)
		for n := len(area) - 1; n >= 0; n-- {
			mask := filter.Get(area[n])
			cellRange[n][1] = max
			if maskMax := mask.Maximum(); maskMax < max {
				max = maskMax
			}
			max--
		}

		for n := 0; n < len(area); n++ {
			if !filter.Restrict(area[n], generator.NewRangeValueMask(cellRange[n][0], cellRange[n][1])) {
				return false
			}
		}
	}
	return true
}

func (r ThermometerRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for _, area := range r.cellAreas[current.Row()][current.Col()] {
		index := 0
		var coordinate grid.Coordinate
		for index, coordinate = range area {
			if coordinate == current {
				break
			}
		}

		max := value - 1
		for n := index - 1; n >= 0; n-- {
			if !state.Restrict(area[n], generator.NewRangeValueMask(1, max)) {
				return
			}
			max--
		}

		min := value + 1
		for n := index + 1; n < len(area); n++ {
			if !state.Restrict(area[n], generator.NewRangeValueMask(min, 9)) {
				return
			}
			min++
		}
	}
	next(state)
}

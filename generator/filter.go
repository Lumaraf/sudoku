package generator

import (
	"github.com/lumaraf/sudoku-checker/grid"
	"math/bits"
)

type Filter struct {
	*ValueMaskGrid
	changed bool
}

func NewFilter(maskGrid *ValueMaskGrid) Filter {
	return Filter{
		maskGrid,
		false,
	}
}

func (f *Filter) Get(coordinate grid.Coordinate) ValueMask {
	return f.ValueMaskGrid[coordinate.Row()][coordinate.Col()]
}

func (f *Filter) Restrict(coordinate grid.Coordinate, mask ValueMask) bool {
	currentMask := &f.ValueMaskGrid[coordinate.Row()][coordinate.Col()]
	initialMask := *currentMask
	*currentMask &= mask
	if *currentMask != initialMask {
		f.changed = true
	}
	return *currentMask != 0
}

func (f *Filter) UniqueGroup(coordinates ...grid.Coordinate) bool {
	solvable := true
	apply := func(filter ValueMask, oneBits int) (changed bool) {
		count := 0
		for _, coordinate := range coordinates {
			mask := &f.ValueMaskGrid[coordinate.Row()][coordinate.Col()]
			if *mask&^filter == 0 {
				count++
			} else if newMask := *mask &^ filter; newMask != *mask {
				*mask = newMask
				changed = true
			}
		}
		if count > oneBits {
			solvable = false
		}
		return
	}

	var scan func(index, count int, mask ValueMask) bool
	scan = func(index, count int, mask ValueMask) bool {
		if index >= 0 {
			mask |= f.ValueMaskGrid.Get(coordinates[index])
			if bits.OnesCount16(uint16(mask)) == count {
				if apply(mask, count) {
					return true
				}
			}
		}
		if count < len(coordinates)-1 {
			index++
			for ; index < len(coordinates); index++ {
				if scan(index, count+1, mask) {
					return true
				}
			}
		}
		return false
	}

	for scan(-1, 0, ValueMask(0)) && solvable {
		f.changed = true
	}

	return solvable
}

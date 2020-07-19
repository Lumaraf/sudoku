package generator

import (
	"github.com/lumaraf/sudoku-checker/grid"
)

type ValueMask uint16

const AllValuesMask = ValueMask(511)

func NewValueMask(values ...uint8) ValueMask {
	m := ValueMask(0)
	for _, value := range values {
		m = m.Set(value)
	}
	return m
}

func (m ValueMask) Get(value uint8) bool {
	return m&(1<<(value-1)) != 0
}

func (m ValueMask) Set(value uint8) ValueMask {
	return m | (1 << (value - 1))
}

func (m ValueMask) Clear(value uint8) ValueMask {
	return m & ^(1 << (value - 1))
}

type ValueMaskGrid [9][9]ValueMask

func NewValueMaskGrid() (mg ValueMaskGrid) {
	for row := 0; row < len(mg); row++ {
		for col := 0; col < len(mg[row]); col++ {
			mg[row][col] = AllValuesMask
		}
	}
	return
}

func (mg *ValueMaskGrid) Get(coordinate grid.Coordinate) ValueMask {
	return mg[coordinate.Row()][coordinate.Col()]
}

func (mg *ValueMaskGrid) CanContain(coordinates []grid.Coordinate, contentMask ValueMask) bool {
	for _, coordinate := range coordinates {
		mask := mg.Get(coordinate)
		if mask&contentMask == 0 {
			return false
		}
	}
	//fmt.Printf("%b %b %b\n", contentMask, combinedMask, contentMask&combinedMask)
	return true
}

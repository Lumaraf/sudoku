package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
)

type KillerEntry struct {
	Area Area
	Sum  uint8
}

func (e KillerEntry) CanSatisfy() bool {
	return true
}

type KillerRule struct {
	entries  []KillerEntry
	entryMap [9][9][]KillerEntry
}

func NewKillerRule(entries []KillerEntry) KillerRule {
	entryMap := [9][9][]KillerEntry{}
	for _, entry := range entries {
		for _, coordinate := range entry.Area {
			entryMap[coordinate.Row()][coordinate.Col()] = append(entryMap[coordinate.Row()][coordinate.Col()], entry)
		}
	}
	return KillerRule{
		entries,
		entryMap,
	}
}

var killerCombinations = generateKillerCombinations(1, 9)

func (r KillerRule) Filter(filter *generator.Filter) bool {
	for _, entry := range r.entries {
		combinations := killerCombinations[len(entry.Area)-1][entry.Sum]
		mask := generator.ValueMask(0)
		for _, combination := range combinations {
			if !filter.CanContain(entry.Area, combination) {
				continue
			}
			mask |= combination
		}
		if mask == 0 {
			return false
		}
		for _, coordinate := range entry.Area {
			if !filter.Restrict(coordinate, mask) {
				return false
			}
		}
	}
	return true
}

func (r KillerRule) Set(current grid.Coordinate, value uint8, state generator.GeneratorState, next generator.NextFunc) {
	for _, entry := range r.entryMap[current.Row()][current.Col()] {
		getAreaMask := func() generator.ValueMask {
			mask := generator.ValueMask(0)
			for _, c := range entry.Area {
				value := state.Get(c)
				if value != 0 {
					mask = mask.Set(value)
				}
			}
			return mask
		}

		areaMask := getAreaMask()
		combinedMask := generator.ValueMask(0)
		combinations := killerCombinations[len(entry.Area)-1][entry.Sum]
		for _, combination := range combinations {
			if combination&areaMask != areaMask {
				continue
			}
			combinedMask |= combination & ^areaMask
		}

		for _, c := range entry.Area {
			if !state.Restrict(c, combinedMask) {
				return
			}
		}
	}
	next(state)
}

func generateKillerCombinations(min, max uint8) (result [9]map[uint8][]generator.ValueMask) {
	for size := 0; size < 9; size++ {
		result[size] = make(map[uint8][]generator.ValueMask)
	}

	var scan func(min uint8, size, sum uint8, mask generator.ValueMask)
	scan = func(min uint8, size, sum uint8, mask generator.ValueMask) {
		for value := min; value <= max; value++ {
			newMask := mask.Set(value)
			result[size-1][sum+value] = append(result[size-1][sum+value], newMask)
			if value < max {
				scan(value+1, size+1, sum+value, newMask)
			}
		}
	}
	scan(min, 1, 0, 0)

	return
}

package rules

import (
	"github.com/lumaraf/sudoku-checker/generator"
	"github.com/lumaraf/sudoku-checker/grid"
)

type MagicSquareRule struct {
	grid.Coordinate
}

var cornerMask = generator.NewValueMask(2, 4, 6, 8)
var sideMask = generator.NewValueMask(1, 3, 7, 9)
var magicSquareRing = []struct {
	c grid.Coordinate
	m generator.ValueMask
}{
	{grid.GetCoordinate(0, 0), cornerMask},
	{grid.GetCoordinate(0, 1), sideMask},
	{grid.GetCoordinate(0, 2), cornerMask},
	{grid.GetCoordinate(1, 2), sideMask},
	{grid.GetCoordinate(2, 2), cornerMask},
	{grid.GetCoordinate(2, 1), sideMask},
	{grid.GetCoordinate(2, 0), cornerMask},
	{grid.GetCoordinate(1, 0), sideMask},
}

func (r MagicSquareRule) Filter(filter *generator.Filter) bool {
	row, col := r.Row(), r.Col()
	if !filter.Restrict(grid.GetCoordinate(row+1, col+1), generator.NewValueMask(5)) {
		return false
	}

	for n, item := range magicSquareRing {
		prev := magicSquareRing[(n+7)%8]
		next := magicSquareRing[(n+1)%8]

		mask := item.m
		mask &= r.getComplementaryMask(filter.Get(r.Coordinate + prev.c))
		mask &= r.getComplementaryMask(filter.Get(r.Coordinate + next.c))
		if !filter.Restrict(r.Coordinate+item.c, mask) {
			return false
		}
	}

	return true
}

func (r MagicSquareRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	for n, item := range magicSquareRing {
		if r.Coordinate+item.c != current {
			continue
		}

		mask := r.getComplementaryMask(generator.NewValueMask(value))
		if prev := magicSquareRing[(n+7)%8]; !state.Restrict(r.Coordinate+prev.c, mask) {
			return
		}
		if next := magicSquareRing[(n+1)%8]; !state.Restrict(r.Coordinate+next.c, mask) {
			return
		}
	}
	next(state)
}

func (r MagicSquareRule) getComplementaryMask(mask generator.ValueMask) (result generator.ValueMask) {
	for n, m := range [9]generator.ValueMask{
		generator.NewValueMask(6, 8),
		generator.NewValueMask(7, 9),
		generator.NewValueMask(4, 8),
		generator.NewValueMask(3, 9),
		generator.NewValueMask(),
		generator.NewValueMask(1, 7),
		generator.NewValueMask(2, 6),
		generator.NewValueMask(1, 3),
		generator.NewValueMask(2, 4),
	} {
		if mask.Get(uint8(n) + 1) {
			result |= m
		}
	}
	return
}

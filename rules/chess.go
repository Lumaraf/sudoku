package rules

import (
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
	"math"
)

type AntiKnightsMoveRule struct {
}

func (r AntiKnightsMoveRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	if current.Row() > 0 {
		if current.Row() > 1 {
			if current.Col() >= 1 && !state.Block(grid.GetCoordinate(current.Row()-2, current.Col()-1), value) {
				return
			}
			if current.Col() <= 7 && !state.Block(grid.GetCoordinate(current.Row()-2, current.Col()+1), value) {
				return
			}
		}
		if current.Col() >= 2 && !state.Block(grid.GetCoordinate(current.Row()-1, current.Col()-2), value) {
			return
		}
		if current.Col() <= 6 && !state.Block(grid.GetCoordinate(current.Row()-1, current.Col()+2), value) {
			return
		}
	}

	if current.Row() < 8 {
		if current.Row() < 7 {
			if current.Col() >= 1 && !state.Block(grid.GetCoordinate(current.Row()+2, current.Col()-1), value) {
				return
			}
			if current.Col() <= 7 && !state.Block(grid.GetCoordinate(current.Row()+2, current.Col()+1), value) {
				return
			}
		}
		if current.Col() >= 2 && !state.Block(grid.GetCoordinate(current.Row()+1, current.Col()-2), value) {
			return
		}
		if current.Col() <= 6 && !state.Block(grid.GetCoordinate(current.Row()+1, current.Col()+2), value) {
			return
		}
	}
	next(state)
}

// it is not possible to place more than 6 groups of queens in a 9x9 (?)
type AntiQueensMoveRule struct {
	mask generator.ValueMask
}

func NewAntiQueensMoveRule(mask generator.ValueMask) AntiQueensMoveRule {
	return AntiQueensMoveRule{
		mask,
	}
}

var queensMoveSets [9][]PlacementSet = generateAntiQueensPlacementSets()

func (r AntiQueensMoveRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	if current.Row() == 0 && r.mask.Get(value) {
		// try out possible placements
		for _, set := range queensMoveSets[current.Col()] {
			set.Place(*state, set, value, next)
		}
	} else {
		next(state)
	}
}

func findQueensMoves(cols ...int) []PlacementSet {
	moveSets := []PlacementSet{}
	queens := PlacementSet{}
	canPlace := func(row, col int) bool {
		for otherRow := 0; otherRow < row; otherRow++ {
			otherCol := queens[otherRow]
			if otherCol == col {
				return false
			}
			if math.Abs(float64(row-otherRow)) == math.Abs(float64(col-otherCol)) {
				return false
			}
		}
		return true
	}
	var trace func(row int)
	trace = func(row int) {
		for col := 0; col < 9; col++ {
			if !canPlace(row, col) {
				continue
			}
			queens[row] = col
			if row == 8 {
				moveSets = append(moveSets, queens)
			} else {
				trace(row + 1)
			}
		}
	}

	for row, col := range cols {
		if !canPlace(row, col) {
			return moveSets
		}
		queens[row] = col
	}
	trace(len(cols))

	return moveSets
}

func generateAntiQueensPlacementSets() (result [9][]PlacementSet) {
	for n := 0; n < 9; n++ {
		result[n] = findQueensMoves(n)
	}
	return
}

type AntiKingsMoveRule struct{}

func (r AntiKingsMoveRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	if current.Row() >= 1 {
		if current.Col() >= 1 && !state.Block(grid.GetCoordinate(current.Row()-1, current.Col()-1), value) {
			return
		}
		if current.Col() <= 7 && !state.Block(grid.GetCoordinate(current.Row()-1, current.Col()+1), value) {
			return
		}
	}

	if current.Row() <= 7 {
		if current.Col() >= 1 && !state.Block(grid.GetCoordinate(current.Row()+1, current.Col()-1), value) {
			return
		}
		if current.Col() <= 7 && !state.Block(grid.GetCoordinate(current.Row()+1, current.Col()+1), value) {
			return
		}
	}
	next(state)
}

type KnightsMoveRule struct {
	mask generator.ValueMask
}

func NewKnightsMoveRule(mask generator.ValueMask) KnightsMoveRule {
	return KnightsMoveRule{
		mask,
	}
}

var knightsMoveSets [9][]PlacementSet = findKnightsPlacementSets()

func (r KnightsMoveRule) Set(current grid.Coordinate, value uint8, state *generator.GeneratorState, next generator.NextFunc) {
	if current.Row() == 0 && r.mask.Get(value) {
		// try out possible placements
		for _, set := range knightsMoveSets[current.Col()] {
			set.Place(*state, set, value, next)
		}
	} else {
		next(state)
	}
}
func findKnightsPlacementSets() (sets [9][]PlacementSet) {
	isHappy := func(set PlacementSet, row int) bool {
		if row < 0 {
			return true
		}

		col := set[row]
		for otherRow := 0; otherRow < row; otherRow++ {
			otherCol := set[otherRow]
			if col == otherCol {
				return false
			}
			if row/3 == otherRow/3 && col/3 == otherCol/3 {
				return false
			}
		}

		if row > 1 && math.Abs(float64(col-set[row-2])) == 1 {
			return true
		}
		if row > 0 && math.Abs(float64(col-set[row-1])) == 2 {
			return true
		}
		if row < 8 && math.Abs(float64(col-set[row+1])) == 2 {
			return true
		}
		if row < 7 && math.Abs(float64(col-set[row+2])) == 1 {
			return true
		}
		return false
	}
	var trace func(set PlacementSet, row int)
	trace = func(set PlacementSet, row int) {
		for col := 0; col < 9; col++ {
			set[row] = col
			if !isHappy(set, row-2) {
				continue
			}
			if row == 8 {
				if isHappy(set, row-1) && isHappy(set, row) {
					sets[set[0]] = append(sets[set[0]], set)
				}
			} else {
				trace(set, row+1)
			}
		}
	}
	trace(PlacementSet{}, 0)

	return
}

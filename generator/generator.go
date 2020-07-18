package generator

import (
	"fmt"
	"github.com/lumaraf/sudoku-checker/grid"
	"math/bits"
)

type Rule interface {
	Set(current grid.Coordinate, value uint8, state *GeneratorState, next NextFunc)
}

type PreMaskingRule interface {
	Rule
	PreMask(maskGrid *ValueMaskGrid)
}

type ValueMask uint16

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

func (mg *ValueMaskGrid) Restrict(coordinate grid.Coordinate, mask ValueMask) {
	mg[coordinate.Row()][coordinate.Col()] &= mask
}

type NextFunc func(state *GeneratorState)

type GeneratorState struct {
	exit    *bool
	steps   *int
	current grid.Coordinate
	values  grid.Grid
	masks   ValueMaskGrid
	rules   []Rule
	result  func(grid.Grid) bool
}

func (g *GeneratorState) Block(coordinate grid.Coordinate, value uint8) bool {
	if currentValue := g.Get(coordinate); currentValue != 0 {
		return true
	}
	m := &g.masks[coordinate.Row()][coordinate.Col()]
	*m = m.Clear(value)
	return *m != 0
}

func (g *GeneratorState) IsAllowed(coordinate grid.Coordinate, value uint8) bool {
	return g.masks[coordinate.Row()][coordinate.Col()].Get(value)
}

func (g *GeneratorState) Get(coordinate grid.Coordinate) uint8 {
	return g.values[coordinate.Row()][coordinate.Col()]
}

func (g *GeneratorState) set(coordinate grid.Coordinate, value uint8) bool {
	cell := &g.values[coordinate.Row()][coordinate.Col()]
	if *cell != 0 {
		return *cell == value
	}
	if !g.IsAllowed(coordinate, value) {
		return false
	}
	*cell = value
	return true
}

func (g GeneratorState) WithCell(coordinate grid.Coordinate, value uint8, next NextFunc) {
	*g.steps++
	if *g.steps%1000000 == 0 {
		// TODO impelemnt progress callbacks
		fmt.Println(*g.steps, "steps")
		g.values.Print()
		fmt.Println()
	}

	if !*g.exit && g.set(coordinate, value) {
		g.processRules(coordinate, value, next)
	}
}

func (g *GeneratorState) processRules(current grid.Coordinate, value uint8, resultCallback NextFunc) {
	rules := g.rules
	var next NextFunc
	next = func(state *GeneratorState) {
		if len(rules) == 0 {
			state.setRestrictedCells(resultCallback)
			return
		}

		rule := rules[0]
		rules = rules[1:]

		rule.Set(current, value, state, next)
	}
	next(g)
}

func (g *GeneratorState) setRestrictedCells(resultCallback NextFunc) {
	coordinates := make([]grid.Coordinate, 0, 81)

	current := g.current
	for current.Row() < 9 && current.Col() < 9 {
		mask := g.masks[current.Row()][current.Col()]
		if bits.OnesCount16(uint16(mask)) == 1 && g.Get(current) == 0 {
			coordinates = append(coordinates, current)
		}

		current++
	}

	var next NextFunc
	next = func(state *GeneratorState) {
		if len(coordinates) == 0 {
			resultCallback(state)
			return
		}
		coordinate := coordinates[0]
		coordinates = coordinates[1:]

		mask := g.masks[coordinate.Row()][coordinate.Col()]
		value := uint8(bits.TrailingZeros16(uint16(mask)) + 1)
		state.WithCell(coordinate, value, next)
	}
	next(g)
}

func (g GeneratorState) FillNext() {
	for g.current.Col() < 9 && g.current.Row() < 9 {
		if g.Get(g.current) == 0 {
			mask := g.masks[g.current.Row()][g.current.Col()]
			value := uint8(0)
			for mask > 0 {
				zeros := uint8(bits.TrailingZeros16(uint16(mask)) + 1)
				mask = mask >> zeros
				value += zeros

				g.WithCell(g.current, value, func(newState *GeneratorState) {
					newState.current++
					newState.FillNext()
				})
			}
			return
		}

		// keep searching for an empty cell
		g.current++
	}

	// no empty cells means the solution is complete
	if !g.result(g.values) {
		*g.exit = true
	}
}

func (g GeneratorState) Dump() {
	g.values.Print()
}

func Generate(rules []Rule, result func(grid.Grid) bool) {
	steps := 0
	exit := false
	g := GeneratorState{
		exit:   &exit,
		steps:  &steps,
		rules:  rules,
		result: result,
	}

	c := grid.Coordinate(0)
	for c.Col() < 9 && c.Row() < 9 {
		g.masks[c.Row()][c.Col()] = 511
		c++
	}

	for _, rule := range rules {
		if preMaskingRule, ok := rule.(PreMaskingRule); ok {
			preMaskingRule.PreMask(&g.masks)
		}
	}

	g.FillNext()
}

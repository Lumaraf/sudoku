package main

import (
	"fmt"
	"github.com/lumaraf/sudoku-checker/generator"
	"github.com/lumaraf/sudoku-checker/grid"
	"github.com/lumaraf/sudoku-checker/rules"
	"time"
)

func main() {
	rules := []generator.Rule{
		&rules.RowRule{},
		&rules.ColumnRule{},
		&rules.SquareRule{},
		rules.NewUniqueAreaRule([]rules.Area{
			{
				grid.GetCoordinate(1, 1), grid.GetCoordinate(1, 2), grid.GetCoordinate(1, 3),
				grid.GetCoordinate(2, 1), grid.GetCoordinate(2, 2), grid.GetCoordinate(2, 3),
				grid.GetCoordinate(3, 1), grid.GetCoordinate(3, 2), grid.GetCoordinate(3, 3),
			},
			{
				grid.GetCoordinate(1, 5), grid.GetCoordinate(1, 6), grid.GetCoordinate(1, 7),
				grid.GetCoordinate(2, 5), grid.GetCoordinate(2, 6), grid.GetCoordinate(2, 7),
				grid.GetCoordinate(3, 5), grid.GetCoordinate(3, 6), grid.GetCoordinate(3, 7),
			},
			{
				grid.GetCoordinate(5, 1), grid.GetCoordinate(5, 2), grid.GetCoordinate(5, 3),
				grid.GetCoordinate(6, 1), grid.GetCoordinate(6, 2), grid.GetCoordinate(6, 3),
				grid.GetCoordinate(7, 1), grid.GetCoordinate(7, 2), grid.GetCoordinate(7, 3),
			},
			{
				grid.GetCoordinate(5, 5), grid.GetCoordinate(5, 6), grid.GetCoordinate(5, 7),
				grid.GetCoordinate(6, 5), grid.GetCoordinate(6, 6), grid.GetCoordinate(6, 7),
				grid.GetCoordinate(7, 5), grid.GetCoordinate(7, 6), grid.GetCoordinate(7, 7),
			},
		}),
		//&rules.AntiKnightsMoveRule{},
		&rules.AntiKingsMoveRule{},
		//rules.NewAntiQueensMoveRule(generator.NewValueMask(1, 2, 3)),
		//rules.NewKnightsMoveRule(generator.NewValueMask(2, 4, 6, 8)),
		&rules.GivenValuesRule{
			grid.GetCoordinate(0, 0): 4,
			grid.GetCoordinate(0, 1): 2,
		},
	}

	start := time.Now()
	count := 0
	generator.Generate(rules, func(g grid.Grid) bool {
		g.Print()
		count++
		return count < 2
	})
	fmt.Println(time.Now().Sub(start))
}

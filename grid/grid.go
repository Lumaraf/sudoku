package grid

import (
	"fmt"
	"strings"
)

type Coordinate uint8

func GetCoordinate(row, col int) Coordinate {
	return Coordinate(row*9 + col)
}

func (c Coordinate) Row() int {
	return int(c) / 9
}

func (c Coordinate) Col() int {
	return int(c) % 9
}

type Grid [9][9]uint8

func (g *Grid) Print() {
	tpl := `╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ % │ % │ % ║ % │ % │ % ║ % │ % │ % ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
`

	sep := strings.Split(tpl, "%")
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			fmt.Print(sep[0])
			sep = sep[1:]
			value := g[row][col]
			if value == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(value)
			}
		}
	}
	fmt.Print(sep[0])
}

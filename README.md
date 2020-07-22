# sudoku

## supported rules

- classic (row, column, square)
- cross
- mirror (linked cells)
- given values
- extra unique areas
- anti knights move
- anti queens move
- anti kings move
- knights move
- non-consecutive
- thermometer
- killer
- magic square

## sudoku definition format

author:
name:
desription:
rules:
  - type: square|row|column|cross|...
  - type: unique-area
    areas:
      - [[1,1],[1,2],[1,3]]

## TODO

- sandwich
- princess move (queen with limited view distance)
- mirror areas
- slow thermos
- constraints between cells and areas (e.g. area == other area or cell1 < cell2)
- import/export rules as yaml

package main

import "fmt"
import "strings"
import "sync"

const (
	alive cell = iota
	dead
	aliveCell = "X"
	deadCell  = "_"
)

func (i cell) String() string {
	switch i {
	case alive:
		return aliveCell
	case dead:
		return deadCell
	}
	return ""
}

type cell int
type row []cell

func (r row) String() string {
	strs := make([]string, len(r))
	for i, cell := range r {
		strs[i] = cell.String()
	}
	return strings.Join(strs, "")
}

// Field is a field of lifegame.
type Field struct {
	height int
	width  int
	rows   []row
}

func newField(height, width int) Field {
	rows := make([]row, height)
	for i := 0; i < height; i++ {
		r := make([]cell, width)
		for j := 0; j < width; j++ {
			r[j] = dead
		}
		rows[i] = r
	}
	return Field{
		height: height,
		width:  width,
		rows:   rows,
	}
}

// CreateField returns a new Field.
func CreateField(input []byte) Field {
	str := strings.TrimSpace(string(input))
	lines := strings.Split(str, "\n")
	h := len(lines)
	w := len(lines[0])
	rows := make([]row, h)
	for i, line := range lines {
		r := make([]cell, w)
		for j := 0; j < w; j++ {
			if line[j:j+1] == aliveCell {
				r[j] = alive
			} else {
				r[j] = dead
			}
		}
		rows[i] = r
	}
	return Field{
		height: h,
		width:  w,
		rows:   rows,
	}
}

// Tick succeeds to next tick.
func (f *Field) Tick() {
	n := tick(*f)
	f.rows = n.rows
}

func tick(f Field) Field {
	outField := newField(f.height, f.width)
	wg := sync.WaitGroup{}
	for i := 0; i < f.height; i++ {
		for j := 0; j < f.width; j++ {
			wg.Add(1)
			go func(x, y int, a, b Field) {
				tickCell(x, y, a, b)
				wg.Done()
			}(i, j, f, outField)
		}
	}
	wg.Wait()
	return outField
}
func (f *Field) setCell(x, y int, c cell) {
	f.rows[y][x] = c
}

func tickCell(x, y int, a, b Field) {
	cond := a.extract(x, y)
	n := cond.nextCell()
	b.setCell(x, y, n)
}

type condField struct {
	*Field
}

func (f *Field) extract(x, y int) condField {
	rows := make([]row, 3)
	for i := range rows {
		yIdx := y - 1 + i
		row := make([]cell, 3)
		for j := range row {
			xIdx := x - 1 + j
			if xIdx < 0 || xIdx >= f.width || yIdx < 0 || yIdx >= f.height {
				row[j] = dead
				continue
			}
			row[j] = f.rows[yIdx][xIdx]
		}
		rows[i] = row
	}
	tmp := &Field{
		width:  3,
		height: 3,
		rows:   rows,
	}
	return condField{tmp}
}

func (c *condField) isAlive() bool {
	return c.rows[1][1] == alive
}

func (c *condField) nextCell() cell {
	/*
		誕生
		死んでいるセルに隣接する生きたセルがちょうど3つあれば、次の世代が誕生する。
		生存
		生きているセルに隣接する生きたセルが2つか3つならば、次の世代でも生存する。
		過疎
		生きているセルに隣接する生きたセルが1つ以下ならば、過疎により死滅する。
		過密
		生きているセルに隣接する生きたセルが4つ以上ならば、過密により死滅する。
	*/
	if !c.isAlive() {
		if c.count() == 3 {
			return alive
		}
		return dead
	}
	switch c.count() {
	case 0, 1:
		return dead
	case 2, 3:
		return alive
	default:
		return dead
	}
}

func (c *condField) count() int {
	sum := 0
	for _, row := range c.rows {
		for _, cell := range row {
			if cell == alive {
				sum++
			}
		}
	}
	if c.isAlive() {
		sum--
	}
	return sum
}

func ShowField(f Field) {
	for _, row := range f.rows {
		fmt.Println(row)
	}
}

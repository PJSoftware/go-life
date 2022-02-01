package cells

import (
	"github.com/PJSoftware/go-life/opengl"
)

type Board struct {
	Grid [][]*Cell
	NumCells int
	gl *opengl.OpenGL
	threshold float32
}

func MakeBoard(numCells int, gl *opengl.OpenGL, threshold float32) Board {
	board := Board{}
	board.Grid = make([][]*Cell, numCells)
	board.NumCells = numCells
	board.gl = gl
	board.threshold = threshold

	for x := 0; x < numCells; x++ {
		board.Grid[x] = make([]*Cell, numCells)
		for y := 0; y < numCells; y++ {
			board.Grid[x][y] = InitCell(x, y, &board)
		}
	}
	return board
}

func (b *Board) UpdateState() {
	for x := range b.Grid {
		for _, c := range b.Grid[x] {
			c.checkState(b)
		}
	}
}

func (b *Board) Draw() {
	b.gl.PreDraw()

	for x := range b.Grid {
		for _, c := range b.Grid[x] {
				c.Draw()
		}
	}

	b.gl.PostDraw()
}

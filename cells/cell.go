package cells

import (
	"math/rand"
)

var (
	unitSquare = []float32{
    -0.5, 0.5, 0,		// triangle, bottom-left
    -0.5, -0.5, 0,
    0.5, -0.5, 0,

    -0.5, 0.5, 0,		// triangle, top-right
    0.5, 0.5, 0,
    0.5, -0.5, 0,
	}
)

type Cell struct {
	drawable uint32

	alive     bool
	aliveNext bool

	x int
	y int

	board *Board
}

func InitCell(x, y int, board *Board) *Cell {
	points := make([]float32, len(unitSquare))
	copy(points, unitSquare)

	spacingPx := float32(board.gl.Size()) / float32(board.NumCells)
	cellSizePx := spacingPx - 2.0
	xPx := spacingPx * (float32(x) + 0.5)
	yPx := spacingPx * (float32(y) + 0.5)
	scaleFactor := 2.0 / float32(board.gl.Size())

	for i := 0; i < len(points); i++ {
		switch i % 3 {
		case 0: // x
			points[i] = (xPx+cellSizePx*points[i])*scaleFactor - 1.0
		case 1: // y
			points[i] = (yPx+cellSizePx*points[i])*scaleFactor - 1.0
		default: // z = 0
			continue
		}
	}

	alive := rand.Float32() < board.threshold

	return &Cell{
		drawable:  board.gl.MakeVao(points),
		alive:     alive,
		aliveNext: alive,
		x:         x,
		y:         y,
		board:     board,
	}
}

func (c *Cell) Draw() {
	if !c.alive {
		return
	}
	c.board.gl.DrawVAO(c.drawable, len(unitSquare))
}

// checkState determines the state of the cell for the next tick of the game.
func (c *Cell) checkState(cells *Board) {
	c.alive = c.aliveNext
	
	liveCount := c.liveNeighbors(cells)
	if c.alive {
		if liveCount < 2 || liveCount > 3 {
			c.aliveNext = false
		}
	} else {
		if liveCount == 3 {
			c.aliveNext = true
		}
	}
}

// liveNeighbors returns the number of live neighbors for a cell.
func (c *Cell) liveNeighbors(cells *Board) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == cells.NumCells {
			x = 0
		} else if x == -1 {
			x = cells.NumCells - 1
		}
		if y == cells.NumCells {
			y = 0
		} else if y == -1 {
			y = cells.NumCells - 1
		}
		if cells.Grid[x][y].alive {
				liveCount++
		}
	}
	
	add(c.x-1, c.y)   // To the left
	add(c.x+1, c.y)   // To the right
	add(c.x, c.y+1)   // up
	add(c.x, c.y-1)   // down
	add(c.x-1, c.y+1) // top-left
	add(c.x+1, c.y+1) // top-right
	add(c.x-1, c.y-1) // bottom-left
	add(c.x+1, c.y-1) // bottom-right
	
	return liveCount
}

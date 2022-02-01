package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/PJSoftware/go-life/cells"
	"github.com/PJSoftware/go-life/opengl"
)

var (
	VERSION = "0.0.2"
	title = "PJSoftware | Conway's Game of Life | v" + VERSION
)

const (
	boardSize = 640 // pixels (square)
	numCells = 64 // cells across and down
	threshold = 0.15 // chance of starting cell being alive
	fps = 12
)

func main() {	
	runtime.LockOSThread()

	glWin := opengl.OpenGLWindow{}
	glWin.Init(boardSize, boardSize, title)

	glShaders := opengl.OpenGLShaders{}
	glShaders.Init("vertexShader", "fragmentShader")

	gl := opengl.OpenGL{}
	gl.Init(&glWin, &glShaders)
	defer gl.Close()
	
	rand.Seed(time.Now().UnixNano())
	board := cells.MakeBoard(numCells, &gl, threshold)
	for gl.IsActive() {
		t := time.Now()
		board.UpdateState()
		board.Draw()
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

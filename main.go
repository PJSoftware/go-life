package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/PJSoftware/go-life/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	VERSION = "0.0.2"

	unitSquare = []float32{
    -0.5, 0.5, 0,		// triangle, bottom-left
    -0.5, -0.5, 0,
    0.5, -0.5, 0,

    -0.5, 0.5, 0,		// triangle, top-right
    0.5, 0.5, 0,
    0.5, -0.5, 0,
	}

	vertexShaderSource = shader.Import("vertexShader")
	fragmentShaderSource = shader.Import("fragmentShader")
)

const (
	boardSize = 640 // pixels (square)
	numCells = 64 // cells across and down
	threshold = 0.15 // chance of starting cell being alive
	fps = 12
)

type cell struct {
	drawable uint32
	
	alive bool
	aliveNext bool

	x int
	y int
}

func main() {	
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	
	program := initOpenGL()
	
	rand.Seed(time.Now().UnixNano())
	cells := makeCells()
	for !window.ShouldClose() {
		t := time.Now()
		for x := range cells {
			for _, c := range cells[x] {
				c.checkState(cells)
			}
	  }
	draw(cells, window, program)
	time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func (c *cell) draw() {
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(unitSquare) / 3))
}

func draw(cells [][]*cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for x := range cells {
		for _, c := range cells[x] {
				c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeCells() [][]*cell {
	cells := make([][]*cell, numCells)
	for x := 0; x < numCells; x++ {
		cells[x] = make([]*cell, numCells)
		for y := 0; y < numCells; y++ {
			cells[x][y] = newCell(x, y)				
		}
	}
	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(unitSquare))
	copy(points, unitSquare)
	
	spacingPx := float32(boardSize) / float32(numCells)
	cellSizePx := spacingPx - 2.0
	xPx := spacingPx * (float32(x) + 0.5)
	yPx := spacingPx * (float32(y) + 0.5)
	scaleFactor := 2.0 / float32(boardSize)

	for i := 0; i < len(points); i++ {
		switch i % 3 {
			case 0: // x
				points[i] = (xPx + cellSizePx * points[i]) * scaleFactor - 1.0
			case 1: // y
				points[i] = (yPx + cellSizePx * points[i]) * scaleFactor - 1.0
			default: // z = 0
				continue
			}
	}

	alive := rand.Float64() < threshold

	return &cell{
		drawable: makeVao(points),
		alive: alive,
		aliveNext: alive,
		x: x,
		y: y,
	}
}

// checkState determines the state of the cell for the next tick of the game.
func (c *cell) checkState(cells [][]*cell) {
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
func (c *cell) liveNeighbors(cells [][]*cell) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == numCells {
			x = 0
		} else if x == -1 {
			x = numCells - 1
		}
		if y == numCells {
			y = 0
		} else if y == -1 {
			y = numCells - 1
		}
		
		if cells[x][y].alive {
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

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(boardSize, boardSize, "PJSoftware | Conway's Game of Life | v" + VERSION, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an initialised program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initialises and returns a vertex array object from the points provided.
func makeVao(points []float32) uint32 {
	var vertexBufferObject uint32
	floatSize := 4	// a float32 takes up 4 bytes in memory
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, floatSize*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	
	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	
	return vertexArrayObject
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	
	cSources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)
	
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	
	return shader, nil
}

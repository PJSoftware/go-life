package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/PJSoftware/go-life/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	square = []float32{
    -0.5, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,

    -0.5, 0.5, 0,
    0.5, 0.5, 0,
    0.5, -0.5, 0,
}
	vertexShaderSource = shader.Import("vertexShader")
	fragmentShaderSource = shader.Import("fragmentShader")
)

const (
	boardDim = 640 // pixels
	boardSize = 28 // cells
)

type cell struct {
	drawable uint32
	
	x int
	y int
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	
	program := initOpenGL()

	cells := makeCells()    
	for !window.ShouldClose() {
			draw(cells, window, program)
	}

}

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square) / 3))
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
	cells := make([][]*cell, boardSize)
	for x := 0; x < boardSize; x++ {
		cells[x] = make([]*cell, boardSize)
			for y := 0; y < boardSize; y++ {
					cells[x][y] = newCell(x, y)
			}
	}
	
	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	copy(points, square)
	
	spacingPx := float32(boardDim) / float32(boardSize)
	cellSizePx := spacingPx - 2.0
	xPx := spacingPx * (float32(x) + 0.5)
	yPx := spacingPx * (float32(y) + 0.5)
	scaleFactor := 2.0 / float32(boardDim)

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

	return &cell{
		drawable: makeVao(points),

		x: x,
		y: y,
	}
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

	window, err := glfw.CreateWindow(boardDim, boardDim, "Conway's Game of Life", nil, nil)
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

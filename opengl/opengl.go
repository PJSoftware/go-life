package opengl

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type OpenGLWindow struct {
	dimX int
	dimY int
	title string
}

func (w *OpenGLWindow) Init(dimX int, dimY int, title string) {
	w.dimX = dimX
	w.dimY = dimY
	w.title = title
}

type OpenGLShaders struct {
	vertexFile string
	fragmentFile string
}

func (s *OpenGLShaders) Init(vert string, frag string) {
	s.vertexFile = vert
	s.fragmentFile = frag
}

type OpenGL struct {
	window *OpenGLWindow
	glProgram uint32
	glWindow  *glfw.Window
}

func (o *OpenGL) Init(win *OpenGLWindow, shader *OpenGLShaders) {
	o.glWindow = initGLWindow(win)
	o.glProgram = initGLProgram(shader)
	o.window = win
}

func (o *OpenGL) Close() {
	glfw.Terminate()
}

func (o *OpenGL) IsActive() bool {
	return !o.glWindow.ShouldClose()
}

func (o *OpenGL) Size() int {
	return o.window.dimX
}

func (o *OpenGL) DrawVAO(vao uint32, size int) {
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(size / 3))
}

func (o *OpenGL) PreDraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(o.glProgram)
}

func (o *OpenGL) PostDraw() {
	glfw.PollEvents()
	o.glWindow.SwapBuffers()
}

// initGLWindow initializes glfw and returns a Window to use.
func initGLWindow(win *OpenGLWindow) *glfw.Window {
	log.Println(fmt.Sprintf("%s (%d x %d)", win.title, win.dimX, win.dimY))
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(win.dimX, win.dimY, win.title, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("glfw.CreateWindow: %s", err.Error()))
	}
	window.MakeContextCurrent()

	return window	
}

// initGLProgram initializes OpenGL and returns an initialised program.
func initGLProgram(shader *OpenGLShaders) uint32 {
	log.Println(fmt.Sprintf("initGLProgram called: %s / %s", shader.fragmentFile, shader.vertexFile))
	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("gl.init: %s", err.Error()))
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(shader.vertexFile, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(shader.fragmentFile, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(glsl string, shaderType uint32) (uint32, error) {
	source := importShader(glsl)
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

func importShader(glslFile string) string {
	content, err := os.ReadFile("shaders/" + glslFile + ".glsl")
	if err != nil {
		log.Fatal(err)
	}
	return string(content) + "\x00"	// shader string must be null-terminated to compile
}

// makeVao initialises and returns a vertex array object from the points provided.
func (o *OpenGL) MakeVao(points []float32) uint32 {
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

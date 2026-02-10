package main

import (
	"fmt"
	"fractlab/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"sync"
	"unsafe"
)

type State struct {
	canvas  *graphics.Canvas
	fractal *graphics.Fractal
	control *ControlState
	wg      *sync.WaitGroup
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func setupWindow() (*glfw.Window, *State) {
	fmt.Print("Initialize GLFW...")
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(mode.Width, mode.Height, "FractLab", nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}

	win.SetCursorPosCallback(cursorPosCallback)
	win.SetScrollCallback(scrollCallback)
	win.SetKeyCallback(keyCallback)

	var state State
	win.SetUserPointer(unsafe.Pointer(&state))
	win.SetPos(0, 0)

	win.MakeContextCurrent()
	glfw.SwapInterval(1) // vsync (set to zero for unlimited framerate
	fmt.Print("Done.\n")

	fmt.Print("Initialize OpenGL...")
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	fmt.Print("Done.\n")

	return win, &state
}

func cleanupWindow(win *glfw.Window) {
	glfw.Terminate()
}

func main() {
	win, state := setupWindow()
	defer cleanupWindow(win)

	width, height := win.GetFramebufferSize()
	canvas := graphics.CanvasNew(mgl32.Vec2{0, 0}, 2, float32(width)/float32(height))
	fractal := graphics.Fractal{
		Colorwheel:      [6]mgl32.Vec3{{0, 0, 1}, {0, 1, 0}, {0, 1, 1}, {1, 0, 0}, {1, 0, 1}, {1, 1, 0}},
		Polynomial:      [3]mgl32.Vec2{{0.325, 0.4}, {0, 0}, {1, 0}},
		DivergenceBound: 2,
		MaxIter:         255,
	}
	c := ControlState{
		MouseX: 0, MouseY: 0,
		Focus:       None,
		Sensitivity: 1.0,
	}
	state.canvas = &canvas
	state.fractal = &fractal
	state.control = &c
	state.wg = new(sync.WaitGroup)

	p, err := graphics.ProgramNew(&canvas, &fractal, win)
	if err != nil {
		panic(err)
	}

	state.control.Scale = state.canvas.Scale
	state.control.Offset = state.control.Offset.Add(mgl32.Vec2{0, 0})

	for !win.ShouldClose() {
		glfw.PollEvents() // only render on change, for continuous drawing use PollEvents instead

		state.canvas.Scale += 0.07 * (state.control.Scale - state.canvas.Scale)
		state.canvas.Offset = state.canvas.Offset.Add(state.control.Offset.Sub(state.canvas.Offset).Mul(0.07))

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		p.Draw()

		win.SwapBuffers()
	}
	state.wg.Wait() // wait for any pending screenshot saving to complete
}

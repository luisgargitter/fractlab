package main

import (
	"fmt"
	"fractlab/fractals"
	"fractlab/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"unsafe"
)

type State struct {
	Animation fractals.Animation
	Viewer    viewerState
	control   ControlState // lowercase not to end up in toml file
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Initialization...")

	v := viewerState{
		OffsetX: 0, OffsetY: 0,
		Scale:   2,
		Fractal: fractals.Mandelbrot(),
	}
	c := ControlState{
		MouseX: 0, MouseY: 0,
		Focus:       None,
		Sensitivity: 0.01,
	}
	a := fractals.Animation{
		Src:  fractals.Mandelbrot(),
		Dest: fractals.Julia(-0.6 + 1i*0.6),
		T:    0.0,
	}

	state := State{
		Animation: a,
		Viewer:    v,
		control:   c,
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	setWindowHints()
	win := initWin()
	setCallbacks(win, unsafe.Pointer(&state))
	width, height := win.GetFramebufferSize()
	state.Viewer.aspectRatio = float32(width) / float32(height)

	program := graphics.BindRenderer(win)
	state.Viewer.uniforms = getUniforms(program)

	canvas := initCanvas()

	for !win.ShouldClose() {
		glfw.PollEvents()

		state.Viewer.Fractal = fractals.GetFractal(state.Animation)
		setUniforms(state.Viewer)

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		canvas.Draw()

		win.SwapBuffers()
	}

}

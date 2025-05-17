package main

import (
	"fractlab/fractals"
	"fractlab/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math"
	"unsafe"
)

type viewerState struct {
	OffsetX, OffsetY, Scale float32
	Fractal                 fractals.Fractal
	aspectRatio             float32
	uniforms                Uniforms
}

type scrollFocus = int

const (
	None = scrollFocus(iota)
	Scale
	Time
)

type ControlState struct {
	MouseX, MouseY float64
	Focus          scrollFocus
	Sensitivity    float32
}

func setWindowHints() {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func initWin() *glfw.Window {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	win, err := glfw.CreateWindow(mode.Width, mode.Height, "FractLab", monitor, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}

	return win
}

func setCallbacks(win *glfw.Window, userPointer unsafe.Pointer) {
	win.SetCursorPosCallback(CursorPosCallback)
	win.SetScrollCallback(ScrollCallback)
	win.SetKeyCallback(KeyCallback)
	win.SetUserPointer(userPointer)
}

func initCanvas() graphics.VAO {
	vertices := []mgl32.Vec3{{-1, -1, 0}, {-1, 3, 0}, {3, -1, 0}}
	surfaces := []graphics.Surface{{0, 1, 2}}
	triangle := graphics.Mesh{Vertices: vertices, Faces: surfaces}
	return triangle.Load()
}

type Uniforms struct {
	Aspect,
	C,
	PZ0,
	PZn,
	OffR,
	OffI,
	Scale int32
}

func getUniforms(program uint32) Uniforms {
	uniforms := Uniforms{}

	uniforms.Aspect = gl.GetUniformLocation(program, gl.Str("aspect\x00"))

	uniforms.C = gl.GetUniformLocation(program, gl.Str("c\x00"))
	uniforms.PZ0 = gl.GetUniformLocation(program, gl.Str("PZ0\x00"))
	uniforms.PZn = gl.GetUniformLocation(program, gl.Str("PZn\x00"))

	uniforms.OffR = gl.GetUniformLocation(program, gl.Str("offR\x00"))
	uniforms.OffI = gl.GetUniformLocation(program, gl.Str("offI\x00"))
	uniforms.Scale = gl.GetUniformLocation(program, gl.Str("scale\x00"))

	return uniforms
}

func setUniforms(v viewerState) {
	f := v.Fractal
	u := v.uniforms
	gl.Uniform1f(u.Aspect, v.aspectRatio)
	gl.Uniform2f(u.C, real(f.C), imag(f.C))
	gl.Uniform2f(u.PZ0, f.PZ0[0], f.PZ0[1])
	gl.Uniform2f(u.PZn, f.PZn[0], f.PZn[1])

	gl.Uniform1f(u.OffR, v.OffsetX)
	gl.Uniform1f(u.OffI, v.OffsetY)
	gl.Uniform1f(u.Scale, v.Scale)
}

func CursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	state := (*State)(w.GetUserPointer())
	width, height := w.GetSize()
	xpos = 2*xpos/float64(width) - 1
	ypos = -(2*ypos/float64(height) - 1)
	deltaMouseX := state.control.MouseX - xpos
	deltaMouseY := state.control.MouseY - ypos
	if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		state.Viewer.OffsetX = state.Viewer.OffsetX + float32(deltaMouseX)*state.Viewer.Scale
		state.Viewer.OffsetY = state.Viewer.OffsetY + float32(deltaMouseY)*state.Viewer.Scale
	}
	state.control.MouseX = xpos
	state.control.MouseY = ypos
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	state := (*State)(w.GetUserPointer())

	s := state.control.Sensitivity * float32(math.Pow(0.95, xoff))
	state.control.Sensitivity = mgl32.Clamp(s, 0, 1)

	switch state.control.Focus {
	case Scale:
		state.Viewer.Scale *= float32(math.Pow(float64(1-state.control.Sensitivity), yoff))
	case Time:
		t := state.Animation.T + float32(yoff)*state.control.Sensitivity
		state.Animation.T = mgl32.Clamp(t, 0, 1)
	}

}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	state := (*State)(w.GetUserPointer())

	switch action {
	case glfw.Press:
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		}

	case glfw.Release:
		switch key {
		case glfw.KeyN:
			state.control.Focus = None
		case glfw.KeyS:
			if mods == glfw.ModControl {
				save(w)
			} else {
				state.control.Focus = Scale
			}
		case glfw.KeyT:
			state.control.Focus = Time
		}

	case glfw.Repeat:

	}
}

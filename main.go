package main

import (
	"fmt"
	"fractlab/graphics"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/BurntSushi/toml"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const degree = 3

type fractalState struct {
	CR, CI float32
	AX, BX float32
	AY, BY float32
}

type viewerState struct {
	OffR, OffI, Scale float32
}

type controlsState struct {
	mouseX, mouseY float64
	zoom           bool
}

type progState struct {
	Fractal  fractalState
	Viewer   viewerState
	controls controlsState // lowercase not to end up in toml file
}

func main() {
	fmt.Println("Initialization...")

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(mode.Width, mode.Height, "FractLab", monitor, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	state := progState{}
	state.Viewer.Scale = 2
	m := fractalState{
		CR: 0,
		CI: 0,
		AX: 0,
		BX: 1,
		AY: 1,
		BY: 0,
	}
	j := fractalState{
		CR: -0.6,
		CI: 0.6,
		AX: 0,
		BX: 0,
		AY: 1,
		BY: 0,
	}
	vecm := mgl32.NewVecNFromData([]float32{m.CR, m.CI, m.AX, m.BX, m.AY, m.BY})
	vecj := mgl32.NewVecNFromData([]float32{j.CR, j.CI, j.AX, j.BX, j.AY, j.BY})

	win.SetUserPointer(unsafe.Pointer(&state))
	win.SetCursorPosCallback(cursorPosCallback)
	win.SetKeyCallback(keyCallback)

	program := graphics.BindRenderer(win)

	uniAspect := gl.GetUniformLocation(program, gl.Str("aspect\x00"))
	width, height := win.GetFramebufferSize()
	gl.Uniform1f(uniAspect, float32(width)/float32(height))

	uniC := gl.GetUniformLocation(program, gl.Str("c\x00"))
	uniPX := gl.GetUniformLocation(program, gl.Str("px\x00"))
	uniPY := gl.GetUniformLocation(program, gl.Str("py\x00"))

	uniOffR := gl.GetUniformLocation(program, gl.Str("offR\x00"))
	uniOffI := gl.GetUniformLocation(program, gl.Str("offI\x00"))
	uniScale := gl.GetUniformLocation(program, gl.Str("scale\x00"))

	vertices := []mgl32.Vec3{{-1, -1, 0}, {-1, 3, 0}, {3, -1, 0}}
	surfaces := []graphics.Surface{{0, 1, 2}}
	triangle := graphics.Mesh{Vertices: vertices, Faces: surfaces}
	canvas := triangle.Load()

	t := float32(0.0)
	for !win.ShouldClose() {
		// not time sensitive
		fractal := state.Fractal
		gl.Uniform2f(uniC, fractal.CR, fractal.CI)
		gl.Uniform2f(uniPX, fractal.AX, fractal.BX)
		gl.Uniform2f(uniPY, fractal.AY, fractal.BY)

		// time sensitive
		glfw.PollEvents()
		gl.Uniform1f(uniOffR, state.Viewer.OffR)
		gl.Uniform1f(uniOffI, state.Viewer.OffI)
		gl.Uniform1f(uniScale, state.Viewer.Scale)
		if state.controls.zoom {
			state.Viewer.Scale = state.Viewer.Scale * 0.997
		}
		if t < 1.0 {
			t += 0.0005
		}
		vect := vecm.Mul(nil, 1-t).Add(nil, vecj.Mul(nil, t))
		fractt := fractalState{
			CR: vect.Get(0),
			CI: vect.Get(1),
			AX: vect.Get(2),
			BX: vect.Get(3),
			AY: vect.Get(4),
			BY: vect.Get(5),
		}
		state.Fractal = fractt

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		canvas.Draw()

		win.SwapBuffers()
	}

}

func save(w *glfw.Window) {
	state := (*progState)(w.GetUserPointer())
	archive(state)
	screenshot(w)
}

func archive(state *progState) {
	file, err := os.Create("saved/" + time.Now().Format(time.DateTime) + ".toml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		log.Fatal(err)
	}
}

func screenshot(w *glfw.Window) {
	width, height := w.GetSize()
	bitmap := make([]uint8, width*height*4)
	gl.Finish() // wait for frame to be done
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(bitmap))
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	copy(img.Pix, bitmap)
	file, err := os.Create("saved/" + time.Now().Format(time.DateTime) + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	state := (*progState)(w.GetUserPointer())
	width, height := w.GetSize()
	xpos = 2*xpos/float64(width) - 1
	ypos = -(2*ypos/float64(height) - 1)
	deltaMouseX := state.controls.mouseX - xpos
	deltaMouseY := state.controls.mouseY - ypos
	if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		state.Viewer.OffR = state.Viewer.OffR + float32(deltaMouseX)*state.Viewer.Scale
		state.Viewer.OffI = state.Viewer.OffI + float32(deltaMouseY)*state.Viewer.Scale
	}
	state.controls.mouseX = xpos
	state.controls.mouseY = ypos
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	state := (*progState)(w.GetUserPointer())

	switch action {
	case glfw.Press:
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		}

	case glfw.Release:
		switch key {
		case glfw.KeyY:
			state.controls.zoom = !state.controls.zoom
		case glfw.KeyR:
			*state = progState{}
			state.Viewer.Scale = 2
		case glfw.KeyS:
			if mods == glfw.ModControl {
				save(w)
			}
		}

	case glfw.Repeat:

	}
}

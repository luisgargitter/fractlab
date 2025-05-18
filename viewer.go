package main

import (
	"fractlab/fractals"
	"fractlab/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

type viewerState struct {
	OffsetX, OffsetY, Scale float32
	Fractal                 fractals.Fractal

	aspectRatio float32
	uniforms    Uniforms
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

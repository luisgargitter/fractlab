package main

import (
	"fmt"
	"fractlab/fractals"
	"fractlab/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type State struct {
	Animation fractals.Animation
	Viewer    viewerState
	control   ControlState // lowercase not to end up in toml file
	wg        sync.WaitGroup
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func fontTest() {
	// Read the font data.
	fontBytes, err := os.ReadFile("fonts/NotoSansMath-Regular.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	rgba := image.NewRGBA(image.Rect(0, 0, 800, 800))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.Pt(0, 0), draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(24)
	c.SetSrc(image.Black)
	c.SetDst(rgba)
	c.SetClip(rgba.Bounds())
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 24)
	if _, err := c.DrawString("\u222B x dx", pt); err != nil {
		log.Println(err)
	}

	file, err := os.Create("saved/" + time.Now().Format(time.UnixDate) + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := png.Encode(file, rgba); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//fontTest()
	fmt.Println("Initialization...")

	v := viewerState{
		OffsetX: 0, OffsetY: 0,
		Scale:   2,
		Fractal: fractals.Mandelbrot(),
	}
	c := ControlState{
		MouseX: 0, MouseY: 0,
		Focus:       None,
		Sensitivity: 1.0,
	}
	a := fractals.Animation{
		Src:  fractals.Mandelbrot(),
		Dest: fractals.Julia(-0.6 + 1i*0.6),
		Time: 0.0,
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

	win := initWin()
	SetCallbacks(win, &state)
	width, height := win.GetFramebufferSize()
	state.Viewer.aspectRatio = float32(width) / float32(height)

	program := graphics.BindRenderer(win)
	state.Viewer.uniforms = getUniforms(program)

	canvas := initCanvas()

	for !win.ShouldClose() {
		glfw.WaitEvents() // only render on change, for continuous drawing use PollEvents instead

		state.Viewer.Fractal = fractals.GetFractal(state.Animation)
		setUniforms(state.Viewer)

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		canvas.Draw()

		win.SwapBuffers()
	}
	state.wg.Wait() // wait for any pending screenshot saving to complete
}

package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Mode = int

const (
	None = Mode(iota)
	Scale
	Time
	Coefficient
	Auto
)

type ControlState struct {
	MouseX, MouseY float64
	Focus          Mode
	Sensitivity    float64
	Scale          float32
	Offset         mgl32.Vec2
}

func normalizeScreenCoordinates(w *glfw.Window, x, y float64) (float64, float64) {
	width, height := w.GetSize()
	return 2*x/float64(width) - 1, -2*y/float64(height) + 1
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	state := (*State)(w.GetUserPointer())

	x, y := normalizeScreenCoordinates(w, xpos, ypos)

	deltaMouseX := float32(state.control.MouseX - x)
	deltaMouseY := float32(state.control.MouseY - y)
	deltaMouse := mgl32.Vec2{deltaMouseX, deltaMouseY}
	state.control.MouseX = x
	state.control.MouseY = y

	if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		state.control.Offset = state.canvas.Offset.Add(deltaMouse.Mul(state.control.Scale))
		state.canvas.Offset = state.control.Offset // for responsiveness
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	state := (*State)(w.GetUserPointer())

	cx, cy := w.GetCursorPos()
	x, y := normalizeScreenCoordinates(w, cx, cy)

	os := state.control.Scale * float32(0.1*yoff)
	state.control.Scale *= float32(1 - 0.1*yoff)
	state.control.Offset = state.control.Offset.Add(mgl32.Vec2{float32(x), float32(y)}.Mul(os))
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	_ = scancode
	state := (*State)(w.GetUserPointer())

	switch action {
	case glfw.Press:
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		}

	case glfw.Release:
		switch key {
		case glfw.KeyS:
			if mods == glfw.ModControl {
				screenshot(w)
			} else {
				state.control.Focus = Scale
			}
		}
	}
}

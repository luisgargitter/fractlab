package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"unsafe"
)

type scrollFocus = int

const (
	None = scrollFocus(iota)
	Scale
	Time
)

type ControlState struct {
	MouseX, MouseY float64
	Focus          scrollFocus
	Sensitivity    float64
}

func SetCallbacks(win *glfw.Window, state *State) {
	win.SetCursorPosCallback(cursorPosCallback)
	win.SetScrollCallback(scrollCallback)
	win.SetKeyCallback(keyCallback)
	win.SetUserPointer(unsafe.Pointer(state))
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	state := (*State)(w.GetUserPointer())
	width, height := w.GetSize()
	xpos = 2*xpos/float64(width) - 1
	ypos = -(2*ypos/float64(height) - 1)
	deltaMouseX := float32(state.control.MouseX - xpos)
	deltaMouseY := float32(state.control.MouseY - ypos)
	if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		state.Viewer.OffsetX = state.Viewer.OffsetX + deltaMouseX*state.Viewer.Scale
		state.Viewer.OffsetY = state.Viewer.OffsetY + deltaMouseY*state.Viewer.Scale
	}
	state.control.MouseX = xpos
	state.control.MouseY = ypos
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	state := (*State)(w.GetUserPointer())

	s := state.control.Sensitivity * math.Pow(0.95, xoff)
	state.control.Sensitivity = mgl64.Clamp(s, 0, 1)

	switch state.control.Focus {
	case Scale:
		state.Viewer.Scale *= float32(math.Pow(1-state.control.Sensitivity*0.1, yoff))
	case Time:
		t := state.Animation.Time + float32(yoff*state.control.Sensitivity*0.01)
		state.Animation.Time = mgl32.Clamp(t, 0, 1)
	}

}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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
				Save(w)
			} else {
				state.control.Focus = Scale
			}
		case glfw.KeyT:
			state.control.Focus = Time

		case glfw.Key1:
			state.control.Sensitivity = math.Pow(0.5, 9)
		case glfw.Key2:
			state.control.Sensitivity = math.Pow(0.5, 8)
		case glfw.Key3:
			state.control.Sensitivity = math.Pow(0.5, 7)
		case glfw.Key4:
			state.control.Sensitivity = math.Pow(0.5, 6)
		case glfw.Key5:
			state.control.Sensitivity = math.Pow(0.5, 5)
		case glfw.Key6:
			state.control.Sensitivity = math.Pow(0.5, 4)
		case glfw.Key7:
			state.control.Sensitivity = math.Pow(0.5, 3)
		case glfw.Key8:
			state.control.Sensitivity = math.Pow(0.5, 2)
		case glfw.Key9:
			state.control.Sensitivity = math.Pow(0.5, 1)
		case glfw.Key0:
			state.control.Sensitivity = math.Pow(0.5, 0)
		}
	case glfw.Repeat:
		switch key {
		case glfw.KeyP:
			if mods == glfw.ModShift { // reverse time
				scrollCallback(w, 0, -1)
			} else {
				scrollCallback(w, 0, 1)
			}

		}

	}
}

package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sqweek/dialog"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func Save(w *glfw.Window) {
	state := (*State)(w.GetUserPointer())
	archive(state)
	screenshot(w)
}

func Load(w *glfw.Window) {
	w.Iconify()
	filename, err := dialog.File().Title("Select a file").Load()
	w.Restore()
	w.Focus()
	if err != nil {
		if !errors.Is(err, dialog.ErrCancelled) {
			log.Fatal(err)
		} else { // file opening aborted
			return
		}
	}
	state := (*State)(w.GetUserPointer())
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// safe current fractal as source for animation
	state.Animation.Src = state.Viewer.Fractal

	decoder := toml.NewDecoder(file)
	if _, err := decoder.Decode(&state.Viewer); err != nil {
		log.Fatal(err)
	}

	// fractal only read to current state not animation
	state.Animation.Dest = state.Viewer.Fractal
	state.Animation.Time = 1.0
}

func archive(state *State) {
	file, err := os.Create("saved/" + time.Now().Format("2006-01-02_15-04-05") + ".toml")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(state.Viewer); err != nil {
		log.Fatal(err)
	}
}

func screenshot(w *glfw.Window) {
	state := (*State)(w.GetUserPointer())

	width, height := w.GetSize()
	bitmap := make([]uint8, width*height*4)
	gl.Finish() // wait for frame to be done
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(bitmap))

	state.wg.Add(1)
	go func() {
		defer state.wg.Done()

		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		copy(img.Pix, bitmap)
		file, err := os.Create("saved/" + time.Now().Format("2006-01-02_15-04-05") + ".png")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		if err := png.Encode(file, img); err != nil {
			log.Fatal(err)
		}
	}()
}

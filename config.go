package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func Save(w *glfw.Window) {
	//state := (*State)(w.GetUserPointer())
	//archive(state)
	screenshot(w)
}

func archive(state *State) {
	file, err := os.Create("saved/" + time.Now().Format(time.UnixDate) + ".toml")
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
	width, height := w.GetSize()
	bitmap := make([]uint8, width*height*4)
	gl.Finish() // wait for frame to be done
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(bitmap))
	go func() {
		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		copy(img.Pix, bitmap)
		file, err := os.Create("saved/" + time.Now().Format(time.UnixDate) + ".png")
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

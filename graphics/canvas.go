package graphics

import "C"
import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Canvas struct {
	Offset mgl32.Vec2
	Scale  float32
	Aspect float32

	vao VAO
}

func (c Canvas) GetSource() string {
	source, err := ShaderSourceFromPath("graphics/shaders/canvas.vert")
	if err != nil {
		panic(err)
	}
	return source
}

func (c Canvas) GetUniformNames() []string {
	return []string{"offset", "scale", "aspect"}
}

func (c Canvas) SetUniforms(u map[string]int32) {
	gl.Uniform2fv(u["offset"], 1, &c.Offset[0])
	gl.Uniform1f(u["scale"], c.Scale)
	gl.Uniform1f(u["aspect"], c.Aspect)
}

func (c Canvas) Draw() {
	c.vao.Draw()
}

func (c *Canvas) load() {
	vertices := []mgl32.Vec3{{-1, -1, 0}, {-1, 1, 0}, {1, 1, 0}, {1, -1, 0}}
	surfaces := []Surface{{0, 1, 2}, {0, 2, 3}}
	triangle := Mesh{Vertices: vertices, Faces: surfaces}
	c.vao = triangle.Load()
}

func CanvasNew(offset mgl32.Vec2, scale float32, aspect float32) Canvas {
	c := Canvas{
		Offset: offset,
		Scale:  scale,
		Aspect: aspect,
	}
	c.load()

	return c
}

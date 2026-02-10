package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Surface [3]uint32

type Mesh struct {
	Vertices []mgl32.Vec3
	Faces    []Surface
}

func (m *Mesh) Load() VAO {
	vbo := constructVBO(m.Vertices)
	ebo := constructEBO(m.Faces)

	return constructVAO(vbo, ebo)
}

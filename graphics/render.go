package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	_ "image/jpeg"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VBO uint32
type EBO struct {
	ebo   uint32
	count int32
}
type VAO struct {
	vao   uint32
	count int32
}

func constructVBO(vertices []mgl32.Vec3) VBO {
	var r uint32
	a := make([][3]float32, len(vertices))
	for i := range vertices {
		v := vertices[i][:]
		a[i] = [3]float32{v[0], v[1], v[2]}
	}

	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(unsafe.Sizeof(a[0]))*len(a),
		unsafe.Pointer(&a[0]),
		gl.STATIC_DRAW,
	)

	return VBO(r)
}

func constructEBO(faces []Surface) EBO {
	var r uint32
	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		int(unsafe.Sizeof(faces[0]))*len(faces),
		unsafe.Pointer(&faces[0]),
		gl.STATIC_DRAW,
	)

	return EBO{r, int32(len(faces))}
}

func constructVAO(vbo VBO, ebo EBO) VAO {
	var r uint32
	gl.GenVertexArrays(1, &r)
	gl.BindVertexArray(r)

	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(vbo))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ebo)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(unsafe.Sizeof([3]float32{})), nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0) // unbind

	return VAO{r, ebo.count}
}

func (v *VAO) Draw() {
	gl.BindVertexArray(v.vao)
	gl.DrawElements(gl.TRIANGLES, int32(v.count*3), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

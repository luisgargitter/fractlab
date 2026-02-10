package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Fractal struct {
	Colorwheel      [6]mgl32.Vec3
	Polynomial      [3]mgl32.Vec2
	DivergenceBound float32
	MaxIter         int32
}

func (v Fractal) GetSource() string {
	source, err := ShaderSourceFromPath("graphics/shaders/fractal.frag")
	if err != nil {
		panic(err)
	}
	return string(source)
}

func (v Fractal) GetUniformNames() []string {
	return []string{"colorwheel", "polynomial", "divergence_bound", "max_iter"}
}

func (v Fractal) SetUniforms(u map[string]int32) {
	gl.Uniform3fv(u["colorwheel"], 6, &v.Colorwheel[0][0])
	gl.Uniform2fv(u["polynomial"], 3, &v.Polynomial[0][0])
	gl.Uniform1f(u["divergence_bound"], v.DivergenceBound)
	gl.Uniform1i(u["max_iter"], v.MaxIter)
}

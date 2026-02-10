package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"os"
	"strings"
)

type Shader interface {
	GetSource() string
	GetUniformNames() []string
	SetUniforms(map[string]int32)
}

type VertexShader interface {
	Shader
	Draw()
}
type FragmentShader Shader

type Program struct {
	id             uint32
	Uniforms       map[string]int32
	vertexShader   VertexShader
	fragmentShader FragmentShader
	window         *glfw.Window
}

func ShaderSourceFromPath(path string) (string, error) {
	source, err := os.ReadFile(path)
	return string(source) + "\x00", err
}

func ProgramNew(vertexShader VertexShader, fragmentShader FragmentShader, window *glfw.Window) (*Program, error) {
	vertexShaderId, err := compileShader(vertexShader.GetSource(), gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShaderId, err := compileShader(fragmentShader.GetSource(), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	p := Program{gl.CreateProgram(),
		make(map[string]int32),
		vertexShader,
		fragmentShader,
		window,
	}

	gl.AttachShader(p.id, vertexShaderId)
	gl.AttachShader(p.id, fragmentShaderId)
	gl.LinkProgram(p.id)

	var status int32
	gl.GetProgramiv(p.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p.id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p.id, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShaderId)
	gl.DeleteShader(fragmentShaderId)

	for _, name := range p.vertexShader.GetUniformNames() {
		p.Uniforms[name] = gl.GetUniformLocation(p.id, gl.Str(name+"\x00"))
	}
	for _, name := range p.fragmentShader.GetUniformNames() {
		p.Uniforms[name] = gl.GetUniformLocation(p.id, gl.Str(name+"\x00"))
	}

	return &p, nil
}

func (p *Program) Draw() {
	p.window.MakeContextCurrent()
	gl.UseProgram(p.id)
	p.vertexShader.SetUniforms(p.Uniforms)
	p.fragmentShader.SetUniforms(p.Uniforms)

	p.vertexShader.Draw()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

#version 330 core

uniform float offR;
uniform float offI;
uniform float scale;
uniform float aspect;

in vec3 vert;

out vec2 coord;

void main() {
    gl_Position = vec4(vert, 1.0f);
    coord = scale * vec2(vert.x * aspect, vert.y) + vec2(offR * aspect, offI);
}
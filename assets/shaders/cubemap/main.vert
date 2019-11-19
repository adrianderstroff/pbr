#version 410 core

layout(location = 0) in vec3 pos;

uniform mat4 M, V, P;

out Vertex {
    vec3 uvw;
} o;

void main() {
    gl_Position = P * V * M * vec4(pos, 1.0);
    o.uvw = pos;
}
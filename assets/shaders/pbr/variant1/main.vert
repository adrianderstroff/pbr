#version 410 core

layout(location = 0) in vec3 pos;
layout(location = 2) in vec3 normal;
layout(location = 1) in vec2 uv;

uniform mat4 M, V, P;

out Vertex {
    vec3 pos;
    vec3 normal;
    vec2 uv;
} o;

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    o.pos    = pos;
    o.normal = normal;
    o.uv     = uv;
}
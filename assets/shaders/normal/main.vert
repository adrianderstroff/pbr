#version 410 core

layout(location = 0) in vec3 pos;
layout(location = 2) in vec3 normal;

uniform mat4 M, V, P;

out Vertex {
    vec3 n;
} o;

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    o.n = normal;
}
#version 430 core

// in
layout(location = 0) in vec3 pos;

// out
out vec3 localPos;

uniform mat4 P, V;

void main() {
    localPos = pos;
    gl_Position = P * V * vec4(pos, 1.0);
}
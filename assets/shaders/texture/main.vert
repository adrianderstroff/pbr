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

#include "constants.glsl"

vec2 spherical(vec3 pos) {
    vec3 dir = normalize(pos);
    float u = 0.5 + atan(dir.x, dir.z) / (2*PI);
    float v = 0.5 + asin(dir.y) / PI;
    return vec2(u, v);
}

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    o.pos    = pos;
    o.normal = normal;
    o.uv     = uv;

    // calculate uv coordinates
    o.uv = spherical(pos);
}
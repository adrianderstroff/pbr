#version 410 core

#include "color.glsl"

in Vertex {
    vec3 n;
} i;

out vec3 fragColor;

void main(){
    fragColor = 0.5*(i.n+1);
}
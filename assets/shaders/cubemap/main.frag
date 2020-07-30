#version 410 core

in Vertex {
    vec3 uvw;
} i;

uniform samplerCube cubemap;

out vec3 color;

#include "tonemapping.glsl"

void main() {             
    color = texture(cubemap, i.uvw).rgb;
    color = tone_mapping(color);
    color = gamma(color);
} 
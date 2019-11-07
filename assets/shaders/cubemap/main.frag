#version 410 core

in Vertex {
    vec3 uvw;
} i;

uniform samplerCube cubemap;

out vec3 color;

void main() {             
    color = vec3(texture(cubemap, i.uvw));
} 
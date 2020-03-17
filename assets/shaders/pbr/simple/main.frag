#version 430 core

//--------------------------------------------------------------------------//
// vertex attributes                                                        //
//--------------------------------------------------------------------------//
in Vertex {
    vec3 pos;
    vec3 normal;
    vec2 uv;
} i;

//--------------------------------------------------------------------------//
// uniforms                                                                 //
//--------------------------------------------------------------------------//
uniform vec3  uCameraPos;
uniform vec3  uLightPos = vec3(20, 20, 20);
uniform vec3  uLightColor = vec3(10);
uniform vec3  uAlbedo;
uniform float uMetallic;
uniform float uRoughness;

//--------------------------------------------------------------------------//
// textures                                                                 //
//--------------------------------------------------------------------------//
layout(binding=0) uniform samplerCube cubemap;
layout(binding=1) uniform sampler2D   noiseTexture;

//--------------------------------------------------------------------------//
// output color                                                             //
//--------------------------------------------------------------------------//
out vec3 outColor;

//--------------------------------------------------------------------------//
// includes                                                                 //
//--------------------------------------------------------------------------//
#include "../shared/util.glsl"
#include "../shared/tonemapping.glsl"
#include "../shared/normal.glsl"
#include "pbr.glsl"

void main(){
    // setup data structures
    PbrMaterial pbr = makePbrMaterial();
    Microfacet micro = makeMicroFacet(pbr, i.pos, normalize(i.pos));

    // determine the reflection depending on the specular properties
    vec3 diffuseColor  = diffuse(pbr);
    vec3 specularColor = specular(pbr, micro);
    vec3 attenuation = PI * specularColor + mix(diffuseColor, vec3(0), pbr.metallic);

    // cosine angle
    float nDotL = max(dot(micro.l, micro.n), 0.0);

    // light intensity
    float dist = length(uLightPos - i.pos);
    vec3 Li = uLightColor / (dist*dist);

    // calculate resulting color
    outColor = attenuation * nDotL * Li;

    outColor = gamma(outColor);
}
#version 430 core

//----------------------------------------------------------------------------//
// vertex attributes                                                          //
//----------------------------------------------------------------------------//
in Vertex {
    vec3 pos;
    vec2 uv;
    vec3 normal;
} i;

//----------------------------------------------------------------------------//
// uniforms                                                                   //
//----------------------------------------------------------------------------//
uniform vec3  uCameraPos;
uniform int   uSamples = 10;
uniform float uGlobalRoughness = 0.1;
uniform float uRandR[100];
uniform float uRandX[100];
uniform float uRandY[100];

//----------------------------------------------------------------------------//
// textures                                                                   //
//----------------------------------------------------------------------------//
layout(binding=0) uniform samplerCube cubemap;
layout(binding=1) uniform sampler2D   albedoTexture;
layout(binding=2) uniform sampler2D   normalTexture;
layout(binding=3) uniform sampler2D   metallicTexture;
layout(binding=4) uniform sampler2D   roughnessTexture;
layout(binding=5) uniform sampler2D   aoTexture;
layout(binding=6) uniform sampler2D   noiseTexture;

//----------------------------------------------------------------------------//
// output color                                                               //
//----------------------------------------------------------------------------//
out vec3 outColor;

//----------------------------------------------------------------------------//
// includes                                                                   //
//----------------------------------------------------------------------------//
#include "../shared/tonemapping.glsl"
#include "pbr.glsl"

void main(){
    // setup parameters
    PbrMaterial pbr   = MakePbrMaterial();
    Microfacet  micro = MakeMicroFacet(pbr, i.pos, i.normal);
    Rand        rand  = MakeRand();

    // calculate for multiple samples
    vec3 color = vec3(0);
    for(int s = 0; s < uSamples; s++) {
        // update random values
        NextRand(rand, s);   

        // calculate resulting color
        vec3 Lo = PI * Brdf(pbr, micro, rand) * CosTheta(micro) * Li(i.pos, micro.l);

        // calculate resulting color
        color += Lo;// * pbr.ao;
    }

    // normalize and map color to LDR then apply gamma function
    vec3 colorHDR = color / uSamples;
    vec3 colorLDR = ReinhardTonemapping(colorHDR);
    outColor      = Gamma(colorLDR);
}
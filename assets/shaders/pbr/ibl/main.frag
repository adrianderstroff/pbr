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
layout (location = 0) out vec3 outColor;
layout (location = 1) out vec3 outAlbedo;
layout (location = 2) out vec3 outNormal;
layout (location = 3) out vec3 outMetallic;
layout (location = 4) out vec3 outRoughness;
layout (location = 5) out vec3 outAo;

//----------------------------------------------------------------------------//
// includes                                                                   //
//----------------------------------------------------------------------------//
#include "../shared/tonemapping.glsl"
#include "pbr.glsl"

void main(){
    // setup parameters
    PbrMaterial pbr   = MakePbrMaterial();
    Microfacet  micro = MakeMicroFacet(pbr, i.pos, i.normal);

    // calculate for multiple samples
    vec3 Lo = vec3(0);
    float dw = 1.0 / uSamples;
    float weights = 0;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.a);
        micro.l = reflect(-micro.v, micro.n);

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        Lo += PI * Brdf2(pbr, micro) * nDotL * Li(i.pos, micro.l) * pbr.ao;
        weights += nDotL;
    }

    // normalize and map color to LDR then apply gamma function
    Lo = Lo / weights;
    vec3 colorLDR = ReinhardTonemapping(Lo);
    outColor      = Gamma(colorLDR);

    // debug
    outAlbedo    = pbr.albedo;
    outNormal    = 0.5 * (1 + micro.n);
    outMetallic  = vec3(pbr.metallic);
    outRoughness = vec3(pbr.roughness);
    outAo        = vec3(pbr.ao);
}
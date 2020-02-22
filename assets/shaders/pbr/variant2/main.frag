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
uniform vec3 uCameraPos;
uniform int  uSamples = 10;
uniform float uGlobalRoughness = 0.1;
uniform float uRandR[100];
uniform float uRandX[100];
uniform float uRandY[100];

//--------------------------------------------------------------------------//
// textures                                                                 //
//--------------------------------------------------------------------------//
layout(binding=0) uniform samplerCube cubemap;
layout(binding=1) uniform sampler2D   albedoTexture;
layout(binding=2) uniform sampler2D   normalTexture;
layout(binding=3) uniform sampler2D   metallicTexture;
layout(binding=4) uniform sampler2D   roughnessTexture;
layout(binding=5) uniform sampler2D   aoTexture;

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
    // grab pbr properties
    PbrMaterial pbr;
    pbr.albedo    = vec3(texture(albedoTexture, i.uv));
    pbr.normal    = vec3(texture(normalTexture, i.uv));
    pbr.metallic  = texture(metallicTexture,    i.uv).x;
    pbr.roughness = texture(roughnessTexture,   i.uv).x;
    pbr.roughness = max(pbr.roughness, uGlobalRoughness);
    pbr.ao        = texture(aoTexture,          i.uv).x;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness;
    pbr.k         = (pbr.a * pbr.a) / 2.0;

    // grab all relevant vectors and the roughness
    Microfacet micro;
    micro.n = normal_mapping(i.normal, pbr.normal);
    micro.v = normalize(uCameraPos - i.pos);

    // setup random variables
    Rand rand;
    rand.r  = 0;
    rand.r1 = 0;
    rand.r2 = 0;

    // calculate for multiple samples
    vec3 color = vec3(0);
    for(int s = 0; s < uSamples; s++) {
        // update random values
        rand.r  = uRandR[s];
        rand.r1 = uRandX[s];
        rand.r2 = uRandY[s];

        // samples the reflected ray using a cosine distribution.
        micro.l = normalize(random_cosine_dir(micro.n, rand.r1, rand.r2, pbr.a));

        // determine ks and kd
        rand.ks = calculateSpecularCoefficient(pbr, micro);
        rand.kd = saturate(1.0 - rand.ks);

        // trace the ray and calculate resulting color
        color += trace(pbr, micro, rand);
    }
    outColor = color / uSamples;
    outColor = tone_mapping(outColor);
}
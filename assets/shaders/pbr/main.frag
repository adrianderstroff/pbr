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
layout(binding=6) uniform sampler2D   noiseTexture;

//--------------------------------------------------------------------------//
// output color                                                             //
//--------------------------------------------------------------------------//
out vec3 outColor;

//--------------------------------------------------------------------------//
// includes                                                                 //
//--------------------------------------------------------------------------//
#include "util.glsl"
#include "pbr.glsl"
#include "tonemapping.glsl"

void main(){
    // grab pbr properties
    PbrMaterial pbr;
    pbr.albedo    = vec3(texture(albedoTexture, i.uv));
    pbr.normal    = vec3(texture(normalTexture, i.uv));
    pbr.metallic  = texture(metallicTexture,    i.uv).x;
    pbr.roughness = texture(roughnessTexture,   i.uv).x;
    pbr.ao        = texture(aoTexture,          i.uv).x;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness;
    pbr.k         = (pbr.a * pbr.a) / 2.0;

    // grab all relevant vectors and the roughness
    Microfacet micro;
    micro.n = normalize(i.normal);
    micro.v = normalize(uCameraPos - i.pos);

    // setup random variables
    Rand rand;
    rand.r  = texture(noiseTexture, i.uv).x;
    rand.r1 = texture(noiseTexture, i.uv).y;
    rand.r2 = texture(noiseTexture, i.uv).z;

    // calculate for multiple samples
    vec3 color = vec3(0);
    vec2 uv = i.uv;
    for(int s = 0; s < uSamples; s++) {
        // update random values
        uv += vec2(0.01, 0.01);
        //uv = vec2(sin(uv.x), cos(uv.y));
        rand.r1 = uRandX[s];
        rand.r2 = uRandY[s];

        // get cosine distributed direction
        micro.l = normalize(random_cosine_dir(micro.n, rand.r1, rand.r2, pbr.a));

        // determine half vector
        micro.h = normalize(micro.l + micro.v);

        // determine ks and kd
        vec3 f = fresnel_schlick(saturate(dot(micro.l, micro.h)), pbr.f0);
        rand.ks = saturate(length(f));
        rand.kd = saturate(1.0 - rand.ks);

        // trace the ray and calculate resulting color
        color += trace(pbr, micro, rand);
    }
    outColor = color / uSamples;
    outColor = tone_mapping(outColor);
}
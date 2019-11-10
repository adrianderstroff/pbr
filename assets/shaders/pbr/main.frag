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
#include "constants.glsl"
#include "util.glsl"
#include "color.glsl"
#include "brdf.glsl"
#include "pbr.glsl"
#include "tonemapping.glsl"

void main(){
    // ray from the camera to the intersection point
    vec3 w0 = normalize(i.pos - uCameraPos);
    vec3 wi = reflect(w0, i.normal);

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
    micro.l = normalize(wi);
    micro.v = normalize(-w0);
    micro.h = normalize(micro.l + micro.v);

    // fresnel property
    vec3 f = fresnel_schlick(saturate(dot(micro.l, micro.h)), pbr.f0);

    // grab random variables
    Rand rand;
    rand.r  = texture(noiseTexture, i.uv).x;
    rand.r1 = texture(noiseTexture, i.uv).y;
    rand.r2 = texture(noiseTexture, i.uv).z;
    rand.kd = saturate(length(f));
    //rand.kd = 1.0;
    rand.ks = saturate(1.0 - rand.kd);

    // calculate for multiple samples
    vec3 color = vec3(0);
    vec2 uv = i.uv;
    for(int s = 0; s < uSamples; s++) {
        // update random values
        uv += vec2(rand.r2, rand.r1);
        rand.r1 = texture(noiseTexture, uv).y;
        rand.r2 = texture(noiseTexture, uv).z;

        // trace the ray and calculate resulting color
        color += trace(w0, wi, pbr, micro, rand);
    }
    outColor = tone_mapping(color / uSamples);
}
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

//----------------------------------------------------------------------------//
// textures                                                                   //
//----------------------------------------------------------------------------//
layout(binding=0) uniform samplerCube cubemap;
layout(binding=1) uniform sampler2D   albedoTexture;
layout(binding=2) uniform sampler2D   normalTexture;
layout(binding=3) uniform sampler2D   metallicTexture;
layout(binding=4) uniform sampler2D   roughnessTexture;
layout(binding=5) uniform sampler2D   aoTexture;

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

vec3 CalculateDiffuseIntegral(PbrMaterial pbr, Microfacet micro) {
    vec3 irradiance = vec3(0);
    float dw = 1.0 / uSamples;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.roughness);
        micro.l = micro.h;

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        irradiance += Li(i.pos, micro.l) * nDotL * dw;
    }
    return diffuse(pbr, micro) * irradiance;
}

vec3 CalculateSpecularIntegral(PbrMaterial pbr, Microfacet micro) {
    vec3 Ls = vec3(0);
    float weights = 0;
    float dw = 1.0 / uSamples;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.roughness);
        micro.l = reflect(-micro.v, micro.h);

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        if (nDotL > 0) {
            Ls += specular(pbr, micro) * nDotL * Li(i.pos, micro.l) * dw;
            weights += nDotL;
        }
    }
    return Ls / weights;
}

vec3 CalculateThemSeparately(PbrMaterial pbr, Microfacet micro) {
    // solve the diffuse integral. we want to calculate the following formula:
    // Ld = kd * c/pi * Int Li(p,wi) * (n . wi) dwi
    vec3 Ld = CalculateDiffuseIntegral(pbr, micro);

    // solve the specular integral. we want to calculate the following formula:
    // Ls = Int fr(p,wi,wo) * Li(p,wi) * (n . wi) dwi
    vec3 Ls = CalculateSpecularIntegral(pbr, micro);

    // add them together
    vec3 Ks = FresnelSchlick(micro.v, micro.n, pbr.f0, pbr.roughness);
    vec3 Kd = (vec3(1) - Ks) * (1-pbr.metallic);
    vec3 Lo = (Kd*Ld + pbr.metallic*Ls) * pbr.ao;

    return Lo;
}

vec3 CalculateThemTogether(PbrMaterial pbr, Microfacet micro) {
    // initialize diffuse and specular term
    vec3 Ld = vec3(0);
    vec3 Ls = vec3(0);

    // integrate the hemisphere using importance sampling
    float weights = 0;
    float dw = 1.0 / uSamples;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.roughness);
        micro.l = reflect(-micro.v, micro.h);

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        if (nDotL > 0) {
            Ls += specular(pbr, micro) * Li(i.pos, micro.l) * nDotL; // * dw
            Ld += Li(i.pos, micro.l) * nDotL * dw;
            weights += nDotL;
        }
    }

    // normalize both terms
    Ls /= weights;
    Ld *= PI * diffuse(pbr, micro);
    
    // determine weights
    vec3 F = FresnelSchlick(micro.v, micro.n, pbr.f0, pbr.roughness);
    vec3 Kd = (vec3(1) - F) * (1-pbr.metallic);
    vec3 Ks = vec3(pbr.metallic);
    
    // put everything together
    vec3 Lo = (Kd*Ld + Ks*Ls) * pbr.ao;
    return Lo;
}

vec3 CalculateBrdfTogether(PbrMaterial pbr, Microfacet micro) {
    // initialize outgoing radiance
    vec3 Lo = vec3(0);

    // integrate the hemisphere using importance sampling
    float weights = 0;
    float dw = 1.0 / uSamples;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.roughness);
        micro.l = reflect(-micro.v, micro.h);

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        if (nDotL > 0) {
            Lo += Brdf(pbr, micro) * Li(i.pos, micro.l) * nDotL * dw;
        }
    }
    
    // put everything together
    Lo *= pbr.ao;
    return Lo;
}

void main(){
    // setup parameters
    PbrMaterial pbr   = MakePbrMaterial();
    Microfacet  micro = MakeMicroFacet(pbr, i.pos, i.normal);

    //vec3 Lo = CalculateThemSeparately(pbr, micro);
    vec3 Lo = CalculateBrdfTogether(pbr, micro);

    // normalize and map color to LDR then apply gamma function
    vec3 colorLDR = ReinhardTonemapping(Lo);
    outColor      = Gamma(colorLDR);

    // debug
    outAlbedo    = pbr.albedo;
    outNormal    = 0.5 * (1 + micro.n);
    outMetallic  = vec3(pbr.metallic);
    outRoughness = vec3(pbr.roughness);
    outAo        = vec3(pbr.ao);
}
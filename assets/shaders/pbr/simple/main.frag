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
// output color                                                             //
//--------------------------------------------------------------------------//
layout (location = 0) out vec3 outColor;
layout (location = 1) out vec3 outDiffuse;
layout (location = 2) out vec3 outSpecular;
layout (location = 3) out vec3 outD;
layout (location = 4) out vec3 outG;
layout (location = 5) out vec3 outF;

//--------------------------------------------------------------------------//
// includes                                                                 //
//--------------------------------------------------------------------------//
#include "../shared/tonemapping.glsl"
#include "../shared/normal.glsl"
#include "pbr.glsl"

void main(){
    // setup data structures
    PbrMaterial pbr = MakePbrMaterial();
    Microfacet micro = MakeMicroFacet(pbr, i.pos, i.normal);

    // cosine angle
    float nDotL = max(dot(micro.l, micro.n), 0.0);

    // calculate resulting color
    vec3 Lo = PI * Brdf(pbr, micro) * nDotL * Li(i.pos);

    // add some ambient lighting
    float ao = 1;
    vec3 ambient = vec3(0.03) * pbr.albedo * ao;
    vec3 colorHDR = ambient + Lo;

    // map HDR to LDR and then map the linear color range to gamma mapped color
    // range.
    vec3 colorLDR = ReinhardTonemapping(colorHDR);
    outColor      = Gamma(colorLDR);

    // calculate for debug purpose
    outDiffuse  = diffuse(pbr);
    outSpecular = specular(pbr, micro);
    outD        = CalcD(pbr, micro);
    outG        = CalcG(pbr, micro);
    outF        = CalcF(pbr, micro);
}
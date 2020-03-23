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
#include "../shared/util.glsl"
#include "../shared/tonemapping.glsl"
#include "../shared/normal.glsl"
#include "pbr.glsl"

vec3 LI() {
    float dist = length(uLightPos - i.pos);
    return uLightColor / (dist*dist);
}

vec3 brdf(in PbrMaterial pbr, in Microfacet micro) {
    vec3 diffuseColor  = mix(diffuse(pbr), vec3(0), uMetallic);
    vec3 specularColor = specular(pbr, micro);
    return specularColor + diffuseColor;
}

void main(){
    // setup data structures
    PbrMaterial pbr = makePbrMaterial();
    Microfacet micro = makeMicroFacet(pbr, i.pos, i.normal);

    // cosine angle
    float nDotL = max(dot(micro.l, micro.n), 0.0);

    // calculate resulting color
    outColor = PI * brdf(pbr, micro) * nDotL * LI();

    // write  resulting colors
    outColor    = gamma(outColor);
    outDiffuse  = diffuse(pbr);
    outSpecular = specular(pbr, micro);
    outD        = calcD(pbr, micro);
    outG        = gamma(calcG(pbr, micro));
    outF        = gamma(calcF(pbr, micro));

    outColor = 0.5*(micro.h + vec3(1));
}
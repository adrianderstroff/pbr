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
uniform vec3 uCameraPos;
uniform int  uSamples = 10;
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
#include "../shared/util.glsl"
#include "../shared/tonemapping.glsl"
#include "../shared/normal.glsl"
#include "pbr.glsl"

void main(){
    // grab pbr properties
    PbrMaterial pbr = makePbrMaterial();

    // grab all relevant vectors and the roughness
    Microfacet micro = makeMicroFacet(pbr, i.pos, i.normal);

    // setup random variables
    Rand rand = makeRand();

    // calculate for multiple samples
    vec3 color = vec3(0);
    for(int s = 0; s < uSamples; s++) {
        // update random values
        nextRand(rand, s);

        // determine ks and kd
        rand.ks = calculateSpecularCoefficient(pbr, micro);
        rand.kd = saturate(1.0 - rand.ks);

        // determine the reflection depending on the specular properties
        vec3 attenuation;
        if(rand.r <= rand.ks) {
            // reflected ray
            micro.l = reflect(-micro.v, micro.n);
            attenuation = rand.ks * specular2(pbr, micro, rand);
        } else {
            // samples the reflected ray using a cosine distribution.
            micro.l = random_cosine_dir2(micro.n, rand.r1, rand.r2, pbr.a);
            attenuation = rand.kd * diffuse(pbr, micro, rand);
        }

        // determine color of indirect light
        vec3 envColor = indirect_light(micro.l, i.pos);

        // calculate resulting color
        color += attenuation * envColor * pbr.ao;
    }
    outColor = color / uSamples;
    outColor = tone_mapping(outColor);
    //outColor = gamma(outColor);
}
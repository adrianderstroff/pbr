#include "constants.glsl"
#include "brdf.glsl"
#include "environment.glsl"

struct PbrMaterial {
    vec3  albedo;
    vec3  normal;
    float metallic;
    float roughness;
    float ao;
    vec3  f0;
    float a;
    float k;
};

struct Microfacet {
    vec3  n;
    vec3  l;
    vec3  v;
    vec3  h;
};

struct Rand {
    float r;
    float r1;
    float r2;
    float kd;
    float ks;
};

vec3 trace(in vec3 w0, in vec3 wi, in PbrMaterial pbr, in Microfacet micro, in Rand rand) {
    vec3 attenuation;

    // get random cosine distributed direction
    micro.h = random_cosine_dir(micro.n, rand.r1, rand.r2, pbr.a);

    if (rand.r <= rand.kd) {
        micro.n = micro.h;
        micro.l = reflect(w0, micro.n);

        float cosine = saturate(dot(micro.n, micro.l));
        float pdf = cosine / PI;

        //attenuation = pbr.albedo * pdf;
        attenuation = pbr.albedo;
    }
    // specular component
    else {
        // compute dot products
        float ndotl = saturate(dot(micro.n, micro.l));
        float ndoth = saturate(dot(micro.n, micro.h));
        float ldoth = saturate(dot(micro.l, micro.h));
        float ndotv = saturate(dot(micro.n, micro.v));

        // calculate brdf
        float d = normal_distribution_ggx(ndoth, pbr.a);
        float g = geometry_smith(ndotl, ndotv, pbr.k);
        vec3  f = fresnel_schlick(ldoth, pbr.f0);
        vec3 ggx = (d * f * g);

        // calculate probability
        float ggxprob = d * ndoth / max(4 * ldoth, 1e-5);
        float pdf = ndotl / max(ggxprob * rand.ks, 1e-5);

        attenuation = ggx * pdf;
    }

    // add ambient occlusion
    attenuation *= pbr.ao;

    // determine indirect illumination
    micro.l = reflect(w0, micro.n);
    vec3 envColor = indirect_light(micro.l, i.pos);

    // calculate resulting color
    return attenuation * envColor;
}
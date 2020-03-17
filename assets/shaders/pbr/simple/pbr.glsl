#include "../shared/constants.glsl"
#include "../shared/brdf.glsl"
#include "../shared/environment.glsl"

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

PbrMaterial makePbrMaterial() {
    PbrMaterial pbr;
    pbr.albedo    = uAlbedo;
    pbr.normal    = vec3(0, 0, 1);
    pbr.metallic  = uMetallic;
    pbr.roughness = uRoughness;
    pbr.ao        = 1.0;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness;
    pbr.k         = ((pbr.a+1) * (pbr.a+1)) / 8.0;
    return pbr;
}

Microfacet makeMicroFacet(in PbrMaterial pbr, vec3 pos, vec3 normal) {
    Microfacet micro;
    micro.n = normal;
    micro.v = normalize(uCameraPos - pos);
    micro.l = normalize(uLightPos - pos);
    return micro;
}

Rand makeRand() {
    Rand rand;
    rand.r  = 0;
    rand.r1 = 0;
    rand.r2 = 0;
    return rand;
}

// calculates the diffuse fraction of the surface.
vec3 diffuse(in PbrMaterial pbr) {
    return pbr.albedo;
}

// calculates the specular fraction of the surface.
vec3 specular(in PbrMaterial pbr, in Microfacet micro) {
    // determine half vector
    micro.h = normalize(micro.l + micro.v);

    // compute dot products
    float ndotl = max(dot(micro.n, micro.l), 0.0);
    float ndoth = dot(micro.n, micro.h);
    float ldoth = dot(micro.l, micro.h);
    float ndotv = dot(micro.n, micro.v);

    // calculate brdf
    float d = normal_distribution_ggx_simple(ndoth, pbr.a);
    float g = geometry_smith_simple(ndotl, ndotv, pbr.k);
    vec3  f = fresnel_schlick(ldoth, pbr.f0);

    return 0.25 * (d * f * g);
}
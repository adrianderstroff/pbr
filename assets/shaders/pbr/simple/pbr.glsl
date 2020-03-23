#include "../shared/constants.glsl"
#include "../shared/brdf.glsl"

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

PbrMaterial makePbrMaterial() {
    PbrMaterial pbr;
    pbr.albedo    = uAlbedo;
    pbr.normal    = i.normal;
    pbr.metallic  = uMetallic;
    pbr.roughness = uRoughness;
    pbr.ao        = 1.0;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness * pbr.roughness;
    pbr.k         = ((pbr.roughness+1) * (pbr.roughness+1)) / 8.0;
    return pbr;
}

Microfacet makeMicroFacet(in PbrMaterial pbr, vec3 pos, vec3 normal) {
    Microfacet micro;
    micro.n = normal;
    micro.v = normalize(uCameraPos - pos);
    micro.l = normalize(uLightPos - pos);
    micro.h = normalize(micro.l + micro.v);
    return micro;
}

// calculates the diffuse fraction of the surface.
vec3 diffuse(in PbrMaterial pbr) {
    return pbr.albedo / PI;
}

// calculates the specular fraction of the surface.
vec3 specular(in PbrMaterial pbr, in Microfacet micro) {
    vec3 v = micro.v;
    vec3 l = micro.l;
    vec3 n = micro.n;
    vec3 h = micro.h;

    // calculate brdf
    float d = normal_distribution_ggx(n, h, pbr.a);
    float g = geometry_smith(l, v, h, pbr.k);
    vec3  f = fresnel_schlick(l, h, pbr.f0);

    // calculate normalization
    float ndotl = max(dot(n, l), 0);
    float ndotv = max(dot(n, v), 0);
    float denom = max(4 * ndotl * ndotv, 0.01);

    return (d * f * g) / denom;
}

// calculates the normal distribution for debugging.
vec3 calcD(in PbrMaterial pbr, in Microfacet micro) {
    vec3 n = micro.n;
    vec3 h = micro.h;
    float d = normal_distribution_ggx(n, h, pbr.a);
    return vec3(d);
}
// calculates the geometric distibution.
vec3 calcG(in PbrMaterial pbr, in Microfacet micro) {
    vec3 v = micro.v;
    vec3 l = micro.l;
    vec3 h = micro.h;
    float g = geometry_smith(l, v, h, pbr.k);
    return vec3(g);
}
// calculates the surface reflectance.
vec3 calcF(in PbrMaterial pbr, in Microfacet micro) {
    vec3 l = micro.l;
    vec3 h = micro.h;
    return fresnel_schlick(l, h, pbr.f0);
}
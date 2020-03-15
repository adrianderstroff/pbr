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
    pbr.albedo    = vec3(texture(albedoTexture, i.uv));
    pbr.albedo    = 3*pbr.albedo;
    pbr.normal    = vec3(texture(normalTexture, i.uv));
    pbr.metallic  = texture(metallicTexture,    i.uv).x;
    pbr.roughness = texture(roughnessTexture,   i.uv).x;
    pbr.roughness = max(pbr.roughness, uGlobalRoughness);
    pbr.ao        = texture(aoTexture,          i.uv).x;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness;
    pbr.k         = (pbr.a * pbr.a) / 2.0;
    return pbr;
}

Microfacet makeMicroFacet(in PbrMaterial pbr, vec3 pos, vec3 normal) {
    Microfacet micro;
    micro.n = normal_mapping(normal, pbr.normal);
    micro.v = normalize(uCameraPos - pos);
    return micro;
}

Rand makeRand() {
    Rand rand;
    rand.r  = 0;
    rand.r1 = 0;
    rand.r2 = 0;
    return rand;
}

/**
 * get the next set of random numbers for the s-th sample.
 */
void nextRand(inout Rand rand, in int s) {
    vec4 rdir = texture(noiseTexture, i.uv);
    rand.r  = fract(uRandR[s] + rdir.z);
    rand.r1 = fract(uRandX[s] + rdir.x);
    rand.r2 = fract(uRandY[s] + rdir.y);
}

/**
 * calculates how reflective the intersected point of the geometry is.
 */
float calculateSpecularCoefficient(in PbrMaterial pbr, in Microfacet micro) {
    // reflected ray
    vec3 l = reflect(-micro.v, micro.n);

    // determine half vector
    vec3 h = normalize(l + micro.v);

    // calculate fresnel
    vec3 f = fresnel_schlick(saturate(dot(l, h)), pbr.f0);

    // calculate specular coefficient
    return saturate(length(f)) * pbr.metallic * pbr.metallic;
    //return saturate(length(f)) * pbr.metallic;
}

float calculateSpecularCoefficient2(in PbrMaterial pbr, in Microfacet micro) {
    return pbr.metallic * pbr.metallic;
}

/**
 * calculates the diffuse fraction of the surface.
 */
vec3 diffuse(in PbrMaterial pbr, in Microfacet micro, in Rand rand) {
    // get pdf
    float cosine = saturate(dot(micro.n, micro.l));
    float pdf = cosine / PI;

    return pbr.albedo * pdf;
}

/**
 * calculates the specular fraction of the surface.
 */
vec3 specular(in PbrMaterial pbr, in Microfacet micro, in Rand rand) {
    // determine half vector
    micro.h = normalize(micro.l + micro.v);

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
    //float ggxprob = d * ndoth / max(4 * ldoth, 1e-5);
    //float pdf = ndotl / max(ggxprob * rand.ks, 1e-5);
    //return ggx * pdf;

    // normalize to get brdf
    vec3 cookTorrance = ggx / max(4 * ndotl * ndotv, 1e-6);

    return cookTorrance;
}

/**
 * calculates the specular fraction of the surface.
 */
vec3 specular2(in PbrMaterial pbr, in Microfacet micro, in Rand rand) {
    // determine half vector
    micro.h = normalize(micro.l + micro.v);

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
    return ggx * pdf;
}
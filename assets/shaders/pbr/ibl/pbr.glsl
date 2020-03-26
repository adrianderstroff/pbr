#include "../shared/constants.glsl"
#include "../shared/brdf.glsl"
#include "../shared/environment.glsl"
#include "../shared/normal.glsl"
#include "../shared/random.glsl"

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

// MakePbrMaterial constructs the PBR Material object
PbrMaterial MakePbrMaterial() {
    PbrMaterial pbr;
    pbr.albedo    = InvGamma(vec3(texture(albedoTexture, i.uv)));
    pbr.normal    = vec3(texture(normalTexture, i.uv));
    pbr.metallic  = texture(metallicTexture,    i.uv).x;
    pbr.roughness = texture(roughnessTexture,   i.uv).x;
    pbr.roughness = max(pbr.roughness, uGlobalRoughness);
    pbr.ao        = texture(aoTexture,          i.uv).x;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness * pbr.roughness;
    pbr.k         = (pbr.a * pbr.a) / 2.0;
    return pbr;
}

// MakeMicroFacet constructs the micro facet object
Microfacet MakeMicroFacet(in PbrMaterial pbr, vec3 pos, vec3 normal) {
    Microfacet micro;
    micro.n = NormalMapping(normal, pbr.normal);
    micro.v = normalize(uCameraPos - pos);
    return micro;
}

// NextRand returns the next set of random numbers for the s-th sample
Rand MakeRand() {
    Rand rand;
    rand.r  = 0;
    rand.r1 = 0;
    rand.r2 = 0;
    return rand;
}

// NextRand returns the next set of random numbers for the s-th sample
void NextRand(inout Rand rand, in int s) {
    vec4 rdir = texture(noiseTexture, i.uv);
    rand.r  = fract(uRandR[s] + rdir.z);
    rand.r1 = fract(uRandX[s] + rdir.x);
    rand.r2 = fract(uRandY[s] + rdir.y);
}

// diffuse calculates the diffuse fraction of the surface.
vec3 diffuse(in PbrMaterial pbr, in Microfacet micro, in Rand rand) {
    // get pdf
    float cosine = Saturate(dot(micro.n, micro.l));
    float pdf = cosine / PI;

    return pbr.albedo * pdf;
}

// specular calculates the specular fraction of the surface
vec3 specular(in PbrMaterial pbr, in Microfacet micro) {
    vec3 v = micro.v;
    vec3 l = micro.l;
    vec3 n = micro.n;
    vec3 h = micro.h;

    // calculate brdf
    float d = NormalDistributionGGX(n, h, pbr.a);
    float g = GeometrySmith(l, v, n, pbr.k);
    vec3  f = FresnelSchlick(v, n, pbr.f0);

    // calculate normalization
    float ndotl = max(dot(n, l), 0);
    float ndotv = max(dot(n, v), 0);
    float denom = max(4 * ndotl * ndotv, 0.01);

    return (d * f * g) / denom;
}

// Li returns the radiance of the light at position pos.
vec3 Li(vec3 pos, vec3 l) {
    return SampleEnvironment(l, pos);
}

float CosTheta(in Microfacet micro) {
    return max(dot(micro.l, micro.n), 0.0);
}

// Brdf calculates the Cook-Torrance BRDF for the given material and surface
// properties.
vec3 Brdf(in PbrMaterial pbr, inout Microfacet micro, in Rand rand) {
    // determine ks and kd
    vec3 Ks = FresnelSchlick(micro.v, micro.n, pbr.f0);
    vec3 Kd = (vec3(1) - Ks) * (1 - pbr.metallic);
    float rs = length(Ks);

    // determine the reflection depending on the specular properties
    vec3 color;
    if(rand.r <= rs) {
        // reflected ray
        micro.l = reflect(-micro.v, micro.n);
        color = specular(pbr, micro);
    } else {
        // samples the reflected ray using a cosine distribution.
        micro.l = RandomCosineDir(micro.n, rand.r1, rand.r2, pbr.a);
        color = Kd * diffuse(pbr, micro, rand);
    }

    return color;
}
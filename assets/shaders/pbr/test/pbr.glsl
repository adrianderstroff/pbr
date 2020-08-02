#include "../shared/constants.glsl"
#include "../shared/brdf.glsl"

struct PbrMaterial {
    vec3  albedo;
    float metallic;
    float roughness;
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

// MakePbrMaterial constructs the PBR Material object
PbrMaterial MakePbrMaterial() {
    PbrMaterial pbr;
    pbr.albedo    = uAlbedo;
    pbr.metallic  = uMetallic;
    pbr.roughness = uRoughness;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness * pbr.roughness;
    pbr.k         = ((pbr.roughness+1) * (pbr.roughness+1)) / 8.0;
    return pbr;
}

// MakeMicroFacet constructs the micro facet object
Microfacet MakeMicroFacet(in PbrMaterial pbr, vec3 pos, vec3 normal) {
    Microfacet micro;
    micro.n = normal;
    micro.v = normalize(uCameraPos - pos);
    micro.l = normalize(uLightPos - pos);
    micro.h = normalize(micro.l + micro.v);
    return micro;
}

// diffuse calculates the diffuse fraction of the surface
vec3 diffuse(in PbrMaterial pbr) {
    return pbr.albedo / PI;
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
vec3 Li(vec3 pos) {
    float dist = length(uLightPos - pos);
    return uLightColor / (dist*dist);
}

// Brdf calculates the Cook-Torrance BRDF for the given material and surface
// properties.
vec3 Brdf(in PbrMaterial pbr, in Microfacet micro) {
    vec3 F = FresnelSchlick(micro.v, micro.n, pbr.f0);
    vec3 kD = (vec3(1) - F) * (1-pbr.metallic);

    vec3 diffuseColor  = mix(diffuse(pbr), vec3(0), uMetallic);
    vec3 specularColor = specular(pbr, micro);

    return specularColor + kD * diffuseColor;
}

// CalcD calculates the normal distribution for debugging.
vec3 CalcD(in PbrMaterial pbr, in Microfacet micro) {
    return vec3(NormalDistributionGGX(micro.n, micro.h, pbr.a));
}
// CalcG calculates the geometric distibution for debugging.
vec3 CalcG(in PbrMaterial pbr, in Microfacet micro) {
    return vec3(GeometrySmith(micro.l, micro.v, micro.n, pbr.k));
}
// CalcF calculates the surface reflectance for debugging.
vec3 CalcF(in PbrMaterial pbr, in Microfacet micro) {
    return FresnelSchlick(micro.v, micro.h, pbr.f0);
}

// GGX version of the geometry function
float chiGGX(float v) {
    return v > 0 ? 1 : 0;
}
vec3 CalcGGGXPartial(in vec3 v, in vec3 n, in vec3 h, in float alpha) {
    float VoH2 = Saturate(dot(v,h));
    float chi = chiGGX(VoH2 / Saturate(dot(v,n)));
    VoH2 = VoH2 * VoH2;
    float tan2 = ( 1 - VoH2 ) / VoH2;
    float val = (chi * 2) / ( 1 + sqrt( 1 + alpha * alpha * tan2 ) );

    return vec3(val);
}
vec3 CalcGGGX(in PbrMaterial pbr, in Microfacet micro) {
    vec3 v = micro.v;
    vec3 l = micro.l;
    vec3 n = micro.n;
    vec3 h = micro.h;
    float alpha = pbr.a;
    
    return CalcGGGXPartial(v, n, h, alpha) * CalcGGGXPartial(l, n, h, alpha);
}
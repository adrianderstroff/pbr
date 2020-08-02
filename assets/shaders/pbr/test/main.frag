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
uniform vec3  uLightPos   = vec3(20, 20, 20);
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
// constants                                                                //
//--------------------------------------------------------------------------//

#define PI 3.1415926535897932384626433832795

//--------------------------------------------------------------------------//
// util functions                                                           //
//--------------------------------------------------------------------------//

// Saturate clamp a value between 0 and 1
float Saturate(float val) {
    return max(0, min(1, val));
}

// Saturate clamps a vector componentwise between 0 and 1
vec3 Saturate(vec3 val) {
    return max(vec3(0), min(vec3(1), val));
}

//--------------------------------------------------------------------------//
// tone mapping                                                             //
//--------------------------------------------------------------------------//

// constants for the tone mapping
const vec3 a = vec3(2.51);
const vec3 b = vec3(0.03);
const vec3 c = vec3(2.43);
const vec3 d = vec3(0.59);
const vec3 e = vec3(0.14);

// Uncharted2Tonemapping performs Uncharted 2's tone mapping
vec3 Uncharted2Tonemapping(in vec3 color) {
	return Saturate((color * (a * color + b)) / (color * (c * color + d) + e));
}

// ReinhardTonemapping performs a simple reinhard tone mapping
vec3 ReinhardTonemapping(in vec3 color) {
	return color / (vec3(1) + color);
}

// Gamma performs gamma mapping 
vec3 Gamma(in vec3 color) {
	return pow(color, vec3(1.0 / 2.2));
}

// InvGamma performs a conversion from sRGB into linear space
float InvGamma(in float val) {
	return pow(val, 2.2);
}

// InvGamma performs a conversion from sRGB into linear space
vec3 InvGamma(in vec3 color) {
	return pow(color, vec3(2.2));
}

//--------------------------------------------------------------------------//
// normal                                                                   //
//--------------------------------------------------------------------------//

// NormalMapping applies normal mapping by taking the relative normal of a 
// normal map and transforming it into the coordinate system defined by the 
// binormal, tangent and normal of the surface. The coordinate system is solely 
// defined by the surface normal of the geometry. the normal map is an RGB 
// texture with the r channel associated with the binormal, the g channel with 
// the tangent and the b channel associated with the surface normal of the 
// coordinate system. both normals are assumed to be length 1.
vec3 NormalMapping(in vec3 surfaceNormal, in vec3 relativeNormal) {
    // calculate tangent and binormal
	vec3 n = normalize(surfaceNormal);
	vec3 t = (dot(n, vec3(0,1,0)) == 0) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

    // in its neutral position the relative normal is pointing in 
    // the z-direction which is blue
    vec3 rn = (2 * relativeNormal) - 1;
    return b * rn.r + t * rn.g + n * rn.b;
}

//--------------------------------------------------------------------------//
// brdf                                                                     //
//--------------------------------------------------------------------------//

// NormalDistributionGGX uses trowbridge-reitz ggx normal distribution function 
// which approximates the relative surface area of microfacets, that are exactly 
// aligned with the halfway vector h. The parameter a specifies the roughness of 
// the surface. Typically a is simply roughness^2.
float NormalDistributionGGX(vec3 n, vec3 h, float a) {
    float angle = max(dot(n, h), 0);
    float a2 = a * a;
    float d = (angle*angle) * (a2-1) + 1;

    return a2 / (PI * d * d);
}

// GeometrySchlickGGX is the schlick approximation of the smith equation. given 
// the surface normal n and a vector v as well as the roughness parameter k, the 
// function approximates how much light can travel in direction v. here k is 
// (roughness+1)^2/8 when doing direct lighting.
float GeometrySchlickGGX(vec3 v, vec3 n, float k) {
    float nDotv = max(dot(n, v), 0);
    return nDotv / (nDotv * (1.0 - k) + k);
}

// GeometrySmith specifies the geometric shadowing of the microfacets based on 
// the view, light and surface normal as well as the roughness of the surface. 
// here k is (roughness+1)^2/8 when doing direct lighting.
float GeometrySmith(vec3 l, vec3 v, vec3 n, float k) {
    float ggx1 = GeometrySchlickGGX(l, n, k);
    float ggx2 = GeometrySchlickGGX(v, n, k);

    return ggx1 * ggx2;
}

// FresnelSchlick specifies the reflection of light on a smooth surface. at a 
// grazing angle all materials become perfect mirrors. f0 is the base 
// reflectivity of the material, which is low for dielectrics (non-metals) and
// usually high for metals. 
vec3 FresnelSchlick(vec3 v, vec3 n, vec3 f0) {
    float vdotn = max(dot(v, n), 0);
    return f0 + (vec3(1) - f0) * pow(1 - vdotn, 5);
}

// FresnelSchlick specifies the reflection of light on a smooth surface. at a 
// grazing angle all materials become perfect mirrors. f0 is the base 
// reflectivity of the material, which is low for dielectrics (non-metals) and
// usually high for metals. 
vec3 FresnelSchlick(vec3 v, vec3 n, vec3 f0, float roughness) {
    float vdotn = max(dot(v, n), 0);
    return f0 + (max(vec3(1 - roughness), f0) - f0) * pow(1 - vdotn, 5);
}

//--------------------------------------------------------------------------//
// pbr                                                                      //
//--------------------------------------------------------------------------//

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

//--------------------------------------------------------------------------//
// main                                                                     //
//--------------------------------------------------------------------------//

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
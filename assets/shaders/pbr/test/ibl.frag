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
uniform int   uSamples = 10;
uniform float uGlobalRoughness = 0.1;

//----------------------------------------------------------------------------//
// textures                                                                   //
//----------------------------------------------------------------------------//
layout(binding=0) uniform samplerCube cubemap;
layout(binding=1) uniform sampler2D   albedoTexture;
layout(binding=2) uniform sampler2D   normalTexture;
layout(binding=3) uniform sampler2D   metallicTexture;
layout(binding=4) uniform sampler2D   roughnessTexture;
layout(binding=5) uniform sampler2D   aoTexture;

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
#define EPS 0.001

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
// environment sampling                                                     //
//--------------------------------------------------------------------------//

// cube dimensions
const vec3 cubeMin    = vec3(-50);
const vec3 cubeMax    = vec3(50);
const vec3 cubeCenter = vec3(0);

// RayBoxIntersection calculates the intersection between a ray and an axis 
// aligned box. 
vec3 RayBoxIntersection(const vec3 boxMin, const vec3 boxMax, const vec3 o, const vec3 dir) {
	vec3 d = (1) * dir;
	
	vec3 tMin = (boxMin - o) / d;
    vec3 tMax = (boxMax - o) / d;
    vec3 t1 = min(tMin, tMax);
    vec3 t2 = max(tMin, tMax);
    float tNear = max(max(t1.x, t1.y), t1.z);
    float tFar = min(min(t2.x, t2.y), t2.z);

	float t = (tNear >= 0) ? tNear : tFar;

	return o + t * d;
}

// SampleEnvironment returns the color of the cubemap of the intersection point
// pos in direction wi.
vec3 SampleEnvironment(in vec3 wi, in vec3 pos) {
    vec3 intersection = RayBoxIntersection(cubeMin, cubeMax, pos, wi);
    vec3 envLookupDir = normalize(intersection - cubeCenter);
    return vec3(texture(cubemap, envLookupDir));
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
	vec3 t = (dot(n, vec3(0,1,0)) < EPS) ? vec3(1, 0, 0) : vec3(0, 1, 0);         // MODIFIED THE CHECK
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

// MakePbrMaterial constructs the PBR Material object
PbrMaterial MakePbrMaterial() {
    PbrMaterial pbr;
    pbr.albedo    = InvGamma(vec3(texture(albedoTexture, i.uv)));
    pbr.normal    = texture(normalTexture, i.uv).xyz;
    pbr.metallic  = texture(metallicTexture,    i.uv).x;
    pbr.roughness = texture(roughnessTexture,   i.uv).x;
    pbr.roughness = max(pbr.roughness, uGlobalRoughness);
    pbr.ao        = texture(aoTexture,          i.uv).x;
    pbr.f0        = mix(vec3(0.04), pbr.albedo, pbr.metallic);
    pbr.a         = pbr.roughness;
    //pbr.a         = pbr.roughness * pbr.roughness;
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

// diffuse calculates the diffuse fraction of the surface.
vec3 diffuse(in PbrMaterial pbr, in Microfacet micro) {
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
    vec3  f = FresnelSchlick(v, n, pbr.f0, pbr.roughness);

    // calculate normalization
    float ndotl = max(dot(n, l), 0);
    float ndotv = max(dot(n, v), 0);
    float denom = max(4 * ndotl * ndotv, 0.01);

    return (d * f * g) / denom;
}

float CosTheta(in Microfacet micro) {
    return max(dot(micro.l, micro.n), 0.0);
}

// Li returns the radiance of the light at position pos.
vec3 Li(vec3 pos, vec3 l) {
    return SampleEnvironment(l, pos);
}

// Brdf calculates the Cook-Torrance BRDF for the given material and surface
// properties.
vec3 Brdf(in PbrMaterial pbr, inout Microfacet micro) {
    // determine ks and kd
    vec3 Ks = FresnelSchlick(micro.v, micro.n, pbr.f0, pbr.roughness);
    vec3 Kd = (vec3(1) - Ks) * (1-pbr.metallic);

    // calculate diffuse and specular term
    vec3 Fs = specular(pbr, micro);
    vec3 Fd = diffuse(pbr, micro);

    return PI * Kd * Fd + pbr.metallic * Fs;
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
    return FresnelSchlick(micro.v, micro.n, pbr.f0);                             // replaced h with n
}

//--------------------------------------------------------------------------//
// random sampling                                                          //
//--------------------------------------------------------------------------//

// radicalInverseVanDerCorpus is an efficient implementation that computes a 
// one-dimensional low discrepancy sequence over the unit interval.
float radicalInverseVanDerCorpus(uint bits) {
	bits = (bits << 16u) | (bits >> 16u);
	bits = ((bits & 0x55555555u) << 1u) | ((bits & 0xAAAAAAAAu) >> 1u);
	bits = ((bits & 0x33333333u) << 2u) | ((bits & 0xCCCCCCCCu) >> 2u);
	bits = ((bits & 0x0F0F0F0Fu) << 4u) | ((bits & 0xF0F0F0F0u) >> 4u);
	bits = ((bits & 0x00FF00FFu) << 8u) | ((bits & 0xFF00FF00u) >> 8u);
	return float(bits) * 2.3283064365386963e-10;
}

// HammersleySampling returns the i-th low discrepancy sample from a set of N
// samples. 
vec2 HammersleySampling(uint i, uint N) {
	return vec2(float(i)/float(N), radicalInverseVanDerCorpus(i));
}

vec3 ImportanceSamplingGGX(vec2 xi, vec3 n, float a) {
	float phi = 2 * PI * xi.x;
	float cosTheta = sqrt((1 - xi.y) / (1 + (a*a - 1) * xi.y));
	float sinTheta = sqrt(1 - cosTheta*cosTheta);

	// spherical to cartesian coordinates
	vec3 pos;
	pos.x = cos(phi) * sinTheta;
	pos.y = sin(phi) * sinTheta;
	pos.z = cosTheta;

	// tangent space to world space
	vec3 up = abs(n.z) < 0.999 ? vec3(0, 0, 1) : vec3(1, 0, 0);
	vec3 tangent = normalize(cross(up, n));
	vec3 bitangent = cross(n, tangent);

	// calculate resulting direction
	return normalize(tangent*pos.x + bitangent*pos.y + n*pos.z);
}

//--------------------------------------------------------------------------//
// main                                                                     //
//--------------------------------------------------------------------------//

vec3 CalculateBrdfTogether(PbrMaterial pbr, Microfacet micro) {
    // initialize outgoing radiance
    vec3 Lo = vec3(0);

    // integrate the hemisphere using importance sampling
    float weights = 0;
    float dw = 1.0 / uSamples;
    for(int s = 0; s < uSamples; s++) {
        // sample direction
        vec2 xi = HammersleySampling(s, uSamples);
        micro.h = ImportanceSamplingGGX(xi, micro.n, pbr.roughness);
        micro.l = reflect(-micro.v, micro.h);

        // calculate angle
        float nDotL = CosTheta(micro);

        // calculate resulting color
        if (nDotL > 0) {
            Lo += Brdf(pbr, micro) * Li(i.pos, micro.l) * nDotL * dw;
        }
    }
    
    // put everything together
    Lo *= pbr.ao;
    return Lo;
}

void main(){
    // renormalize normal after rasterization
    vec3 n = normalize(i.normal);                                                // added to fix grid artifacts

    // setup parameters
    PbrMaterial pbr   = MakePbrMaterial();
    Microfacet  micro = MakeMicroFacet(pbr, i.pos, n);

    //vec3 Lo = CalculateThemSeparately(pbr, micro);
    vec3 Lo = CalculateBrdfTogether(pbr, micro);

    // normalize and map color to LDR then apply gamma function
    vec3 colorLDR = ReinhardTonemapping(Lo);
    outColor      = Gamma(colorLDR);

    // calculate for debug purpose
    outDiffuse  = Gamma(diffuse(pbr, micro));
    outSpecular = Gamma(specular(pbr, micro));
    outD        = CalcD(pbr, micro);
    outG        = CalcG(pbr, micro);
    outF        = CalcF(pbr, micro);
}
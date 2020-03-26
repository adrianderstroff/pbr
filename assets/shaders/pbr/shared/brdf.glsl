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
/**
 * trowbridge-reitz ggx normal distribution function
 * approximates the relative surface area of microfacets, that are
 * exactly aligned with the halfway vector h. The parameter a 
 * specifies the roughness of the surface.
 */
float normal_distribution_ggx(float ndoth, float a) {
    float a2 = a * a;
    float d = (ndoth * (a2*ndoth - ndoth) + 1.0);

    return a2 / (PI * d * d);
}

float normal_distribution_ggx_simple(float ndoth, float a) {
    float a2 = a * a * a * a;
    float d = (ndoth * ndoth) * (a2 - 1) + 1;

    return a2 / (PI * d * d);
}

// schlick approximation of the smith equation
float geometry_schlick_ggx(float ndotv, float k) {
    return ndotv / (ndotv * (1.0 - k) + k);
}

float geometry_schlick_ggx_simple(float ndotv, float k) {
    return 1.0 / (ndotv * (1.0 - k) + k);
}

/**
 * specifies the geometric shadowing of the microfacets based
 * on the view vector and the roughness of the surface. here
 * k depends on the roughness of the surface. this implementation
 * uses the smith method.
 */
float geometry_smith(float ndotl, float ndotv, float k) {
    float ggx1 = geometry_schlick_ggx(ndotv, k);
    float ggx2 = geometry_schlick_ggx(ndotl, k);

    return ggx1 * ggx2;
}

float geometry_smith_simple(float ndotl, float ndotv, float k) {
    float ggx1 = geometry_schlick_ggx_simple(ndotv, k);
    float ggx2 = geometry_schlick_ggx_simple(ndotl, k);

    return ggx1 * ggx2;
}

// specifies the reflection of light on a smooth surface
vec3 fresnel_schlick(float ldoth, vec3 f0) {
    return f0 + (vec3(1.0) - f0) * pow(1.0 - ldoth, 5);
}
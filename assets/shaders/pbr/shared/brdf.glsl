/**
 * trowbridge-reitz ggx normal distribution function approximates the relative 
 * surface area of microfacets, that are exactly aligned with the halfway 
 * vector h. The parameter a specifies the roughness of the surface.
 */
float normal_distribution_ggx(vec3 n, vec3 h, float a) {
    float angle = max(dot(n, h), 0);
    float a2 = a * a;
    float d = (angle*angle) * (a2-1) + 1;

    return a2 / (PI * d * d);
}

// schlick approximation of the smith equation
float geometry_schlick_ggx(vec3 v, vec3 h, float k) {
    float angle = max(dot(h, v), 0);
    return angle / (angle * (1.0 - k) + k);
}

/**
 * specifies the geometric shadowing of the microfacets based on the view vector 
 * and the roughness of the surface. here k depends on the roughness of the 
 * surface. this implementation uses the smith method.
 */
float geometry_smith(vec3 l, vec3 v, vec3 h, float k) {
    float ggx1 = geometry_schlick_ggx(l, h, k);
    float ggx2 = geometry_schlick_ggx(v, h, k);

    return ggx1 * ggx2;
}

// specifies the reflection of light on a smooth surface
vec3 fresnel_schlick(vec3 l, vec3 h, vec3 f0) {
    float ldoth = max(dot(l, h), 0);
    return f0 + (vec3(1) - f0) * pow(1 - ldoth, 5);
}
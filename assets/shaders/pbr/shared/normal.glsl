#include "constants.glsl"

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
	vec3 t = (dot(n, vec3(0,1,0)) < EPS) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

    // in its neutral position the relative normal is pointing in 
    // the z-direction which is blue
    vec3 rn = (2 * relativeNormal) - 1;
    return b * rn.r + t * rn.g + n * rn.b;
}
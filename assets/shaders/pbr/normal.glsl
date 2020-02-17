vec3 normal_mapping(in vec3 surfaceNormal, in vec3 relativeNormal) {
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
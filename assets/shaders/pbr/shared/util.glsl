#include "constants.glsl"

/**
 * util function to clamp a value between 0 and 1
 */
float saturate(float val) {
    return max(0, min(1, val));
}

/**
 * util function to clamp a vector componentwise between 0 and 1
 */
vec3 saturate(vec3 val) {
    return max(vec3(0), min(vec3(1), val));
}

/**
 * returns a random normal following a cosine distribution
 */
vec3 random_cosine_dir(in vec3 normal, float r1, float r2, float a) {
	// calculate tangent and binormal
	vec3 n = normalize(normal);
	vec3 t = (dot(n, vec3(0,1,0)) == 0) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

	float phi = 2 * PI * r1;

	float x = cos(phi) * sqrt(r2) * a;
	float y = sqrt(1 - r2);
	float z = sin(phi) * sqrt(r2) * a;

	vec3 dir = (t * x) + (n * y) + (b * z);
	return normalize(dir);
}

/**
 * returns a random normal following a cosine distribution
 */
vec3 random_cosine_dir2(in vec3 normal, float r1, float r2, float a) {
	// calculate tangent and binormal
	vec3 n = normalize(normal);
	vec3 t = (dot(n, vec3(0,1,0)) == 0) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

	float sinTheta = sqrt(r1);
	float cosTheta = sqrt(1 - sinTheta*sinTheta);
	float psi = 2 * PI * r2;

	vec3 v1 = cosTheta * n;
	vec3 v2 = sinTheta * cos(psi) * t;
	vec3 v3 = sinTheta * sin(psi) * b;

	vec3 dir = v1 + v2 + v3;
	return normalize(dir);
}

/**
 * calculates the intersection between a ray and an axis aligned box
 */
vec3 ray_box_intersection(const vec3 boxMin, const vec3 boxMax, const vec3 o, const vec3 dir) {
	vec3 d = (-1) * dir;
	
	vec3 tMin = (boxMin - o) / d;
    vec3 tMax = (boxMax - o) / d;
    vec3 t1 = min(tMin, tMax);
    vec3 t2 = max(tMin, tMax);
    float tNear = max(max(t1.x, t1.y), t1.z);
    float tFar = min(min(t2.x, t2.y), t2.z);

	return o + tNear * d;
}
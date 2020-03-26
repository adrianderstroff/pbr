#include "constants.glsl"

// RandomCosineDir returns a random normal following a cosine distribution. the
// main direction of the cosine distribution is specified by the normal. r1 and
// r2 are two random numbers that are used to calculate the resulting direction
// while a is the surface roughness that turns the half sphere into an ellipsoid
vec3 RandomCosineDir(in vec3 normal, float r1, float r2, float a) {
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

// RandomCosineDir returns a random normal following a cosine distribution. the
// main direction of the cosine distribution is specified by the normal. r1 and
// r2 are two random numbers that are used to calculate the resulting direction
// while a is the surface roughness that turns the half sphere into an ellipsoid
// https://www.particleincelsl.com/2015/cosine-distribution/
vec3 RandomCosineDir2(in vec3 normal, float r1, float r2, float a) {
	// calculate tangent and binormal
	vec3 n = normalize(normal);
	vec3 t = (dot(n, vec3(0,1,0)) == 1) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

	float sinTheta = sqrt(r1);
	float cosTheta = sqrt(1 - sinTheta*sinTheta);
	float psi = 2 * PI * r2;

	vec3 v1 = cosTheta * n;
	vec3 v2 = sinTheta * cos(psi) * t * a;
	vec3 v3 = sinTheta * sin(psi) * b * a;

	vec3 dir = v1 + v2 + v3;
	return normalize(dir);
}
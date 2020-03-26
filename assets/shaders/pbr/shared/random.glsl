#include "constants.glsl"

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
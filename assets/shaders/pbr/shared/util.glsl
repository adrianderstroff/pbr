#include "constants.glsl"
#include "color.glsl"

// calculates the color from the direction projected onto the plane defined
// by the normal.
vec3 ColorDirection(in vec3 normal, in vec3 dir) {
	// calculate tangent and binormal
	vec3 n = normalize(normal);
	vec3 t = (dot(n, vec3(0,1,0)) == 1) ? vec3(1, 0, 0) : vec3(0, 1, 0);
	vec3 b = cross(t, n);
	t = cross(n, b);

	float x = dot(dir, b);
	float y = dot(dir, t);

	float h = (atan(-y, -x) + PI) / (2*PI);
	float s = 1;
	float v = length(vec2(x, y));
	v = 1;

	return Hsv2Rgb(vec3(h, s, v));
}
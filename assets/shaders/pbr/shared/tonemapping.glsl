#include "math.glsl"

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
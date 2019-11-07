// constants for the tone mapping
const vec3 a = vec3(2.51);
const vec3 b = vec3(0.03);
const vec3 c = vec3(2.43);
const vec3 d = vec3(0.59);
const vec3 e = vec3(0.14);

vec3 tone_mapping(vec3 color) {
	return saturate((color * (a * color + b)) / (color * (c * color + d) + e));
}
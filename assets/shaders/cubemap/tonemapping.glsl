/**
 * util function to clamp a vector componentwise between 0 and 1
 */
vec3 saturate(vec3 val) {
    return max(vec3(0), min(vec3(1), val));
}

// constants for the tone mapping
const vec3 a = vec3(2.51);
const vec3 b = vec3(0.03);
const vec3 c = vec3(2.43);
const vec3 d = vec3(0.59);
const vec3 e = vec3(0.14);

/**
 * Uncharted 2 tone mapping
 */
vec3 tone_mapping(in vec3 color) {
	return saturate((color * (a * color + b)) / (color * (c * color + d) + e));
}

/**
 * simple reinhard 
 */
vec3 reinhard_tonemapping(in vec3 color) {
	return color / (vec3(1) + color);
}

/**
 * simple reinhard 
 */
vec3 gamma(in vec3 color) {
	return pow(color, vec3(1.0 / 2.2));
}
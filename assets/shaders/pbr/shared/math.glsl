// Saturate clamp a value between 0 and 1
float Saturate(float val) {
    return max(0, min(1, val));
}

// Saturate clamps a vector componentwise between 0 and 1
vec3 Saturate(vec3 val) {
    return max(vec3(0), min(vec3(1), val));
}
// cube dimensions
const vec3 cubeMin = vec3(-30);
const vec3 cubeMax = vec3(30);
const vec3 cubeCenter = vec3(0);

/**
 * determines the light color of the surrounding environment
 */
vec3 indirect_light(in vec3 wi, in vec3 pos) {
    vec3 intersection = ray_box_intersection(cubeMin, cubeMax, pos, wi);
    vec3 envLookupDir = normalize(intersection - cubeCenter);
    return vec3(texture(cubemap, envLookupDir));
}
// cube dimensions
const vec3 cubeMin = vec3(-50);
const vec3 cubeMax = vec3(50);
const vec3 cubeCenter = vec3(0);

// RayBoxIntersection calculates the intersection between a ray and an axis 
// aligned box. 
vec3 RayBoxIntersection(const vec3 boxMin, const vec3 boxMax, const vec3 o, const vec3 dir) {
	vec3 d = (1) * dir;
	
	vec3 tMin = (boxMin - o) / d;
    vec3 tMax = (boxMax - o) / d;
    vec3 t1 = min(tMin, tMax);
    vec3 t2 = max(tMin, tMax);
    float tNear = max(max(t1.x, t1.y), t1.z);
    float tFar = min(min(t2.x, t2.y), t2.z);

	float t = (tNear >= 0) ? tNear : tFar;

	return o + t * d;
}

// SampleEnvironment returns the color of the cubemap of the intersection point
// pos in direction wi.
vec3 SampleEnvironment(in vec3 wi, in vec3 pos) {
    vec3 intersection = RayBoxIntersection(cubeMin, cubeMax, pos, wi);
    vec3 envLookupDir = normalize(intersection - cubeCenter);
    return vec3(texture(cubemap, envLookupDir));
}
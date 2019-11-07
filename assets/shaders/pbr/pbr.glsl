struct PbrMaterial {
    vec3  albedo;
    vec3  normal;
    float metallic;
    float roughness;
    float ao;
    vec3  f0;
    float a;
    float k;
};

struct Microfacet {
    vec3  n;
    vec3  l;
    vec3  v;
    vec3  h;
};

struct Rand {
    float r;
    float r1;
    float r2;
    float kd;
    float ks;
};

vec3 trace(vec3 w0, vec3 wi, PbrMaterial pbr, in Microfacet micro, Rand rand) {
    vec3 attenuation;

    if (rand.r <= rand.kd) {
        // get random cosine distributed direction
        micro.l = random_cosine_dir(micro.n, rand.r1, rand.r2, 0.1);
        //micro.l = micro.n;

        float cosine = saturate(dot(micro.n, micro.l));
        float pdf = cosine / PI;

        //attenuation = pbr.albedo * pdf;
        attenuation = pbr.albedo;
    }
    // specular component
    else {
        // compute dot products
        float ndotl = saturate(dot(micro.n, micro.l));
        float ndoth = saturate(dot(micro.n, micro.h));
        float ldoth = saturate(dot(micro.l, micro.h));
        float ndotv = saturate(dot(micro.n, micro.v));

        // calculate brdf
        float d = normal_distribution_ggx(ndoth, pbr.a);
        float g = geometry_smith(ndotl, ndotv, pbr.k);
        vec3  f = fresnel_schlick(ldoth, pbr.f0);
        vec3 ggx = (d * f * g);

        // calculate probability
        float ggxprob = d * ndoth / max(4 * ldoth, 1e-6);
        float pdf = ndotl / max(ggxprob * rand.ks, 1e-6);

        attenuation = ggx * pdf;
    }

    // add ambient occlusion
    //attenuation *= pbr.ao;

    // determine indirect illumination
    vec3 envLookupDir = determine_direction(micro.l, vec3(0,0,0), i.pos);
    vec3 envColor = vec3(texture(cubemap, envLookupDir));
    //envColor = 0.5*(normalize(micro.l) + vec3(1));

    // calculate resulting color
    return attenuation * envColor;
    //return envColor;
}
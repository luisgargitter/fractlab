#version 330 core

uniform vec3[6] colorwheel;

uniform vec2[3] polynomial;
uniform float divergence_bound;
uniform int max_iter;

in vec2 coord;

out vec4 color;

vec2 cMul(vec2 a, vec2 b) {
    return vec2(
    a.x * b.x - a.y * b.y,  // real part
    a.x * b.y + a.y * b.x   // imag part
    );
}

vec2 evaluate_polynomial(vec2 z) {
    vec2 a = polynomial[0];
    vec2 zp = z;
    for(int i = 1; i < 3; i++) {
        a = a + cMul(polynomial[i], zp);
        zp = cMul(zp, z);
    }
    return a;
}

// calculates divergence bound for this interpolated fractal (analytically derived)
float divBound(float alpha, float beta, vec2 prec) {
    alpha = abs(alpha);
    beta = abs(beta);
    float co = length(prec);
    float p2 = (beta+1)/(2*alpha);
    float q = -co/alpha;
    float R = p2 + sqrt(p2*p2 - q);

    return R;
}

vec3 colorFromHueSat(float hue, float sat) {
    float h6 = hue*6.0;
    int i = int(h6);
    vec3 c1 = colorwheel[i];
    vec3 c2 = colorwheel[(i+1) % 6];
    float t = h6 - float(i);
    vec3 col = c1 * (1-t) + c2 * t;
    return sat * col/length(col);
}

vec4 fixpoint_iteration(vec2 z) {
    int i;
    for(i = 0; (i < max_iter) && (length(z) < divergence_bound); i++) {
        z = evaluate_polynomial(z); // general julia set
    }
    float hue = float(i)/float(max_iter);
    float sat = 1 - exp(-float(max_iter-i));

    return vec4(colorFromHueSat(hue, sat), 1.0);
}

void main() {
    color = fixpoint_iteration(coord);
}

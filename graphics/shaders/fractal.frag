#version 330 core

uniform int overlay;
uniform float scale; // for ui elements

uniform vec2 c;
uniform vec2 PZ0;
uniform vec2 PZn;

in vec2 coord;

out vec4 color;

vec2 cMul(vec2 a, vec2 b) {
    return vec2(
    a.x * b.x - a.y * b.y,  // real part
    a.x * b.y + a.y * b.x   // imag part
    );
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

const vec3 Hue[6] = vec3[](vec3(0, 0, 1), vec3(0, 1, 0), vec3(0, 1, 1), vec3(1, 0, 0), vec3(1, 0, 1), vec3(1, 1, 0));
vec3 colorFromHueSat(float hue, float sat) {
    float h6 = hue*6.0;
    int i = int(h6);
    vec3 c1 = Hue[i];
    vec3 c2 = Hue[(i+1) % 6];
    float t = h6 - float(i);
    vec3 col = c1 * (1-t) + c2 * t;
    return sat * col/length(col);
}

const int depth = 256;

void main()
{
    vec2 x = coord;
    vec2 x2 = cMul(x, x);
    vec2 prec = PZ0[0]*x2 + PZ0[1]*x + c;
    float R = divBound(PZn[0], PZn[1], prec);

    vec2 y = x;

    int j = 0;

    for(int i = 0; i < depth && length(y) < R; i++) {
        y = PZn[0]*cMul(y, y) + PZn[1]*y + prec;

        j++;
    }
    float hue = float(j)/float(depth);
    float sat = 1 - pow(0.5, float(depth-j));

    color = vec4(colorFromHueSat(hue, sat), 1.0);

    if(overlay == 1) {
        float lx = length(x);
        bool circle = abs(1 - length(x)) < 0.001 * scale;
        bool xaxis = abs(0 - x.x) < 0.001 * scale;
        bool yaxis = abs(0 - x.y) < 0.001 * scale;
        bool coeff = abs(length(x - c)) < 0.007 * scale;
        if (circle || xaxis || yaxis || coeff) {
            color = vec4(1, 1, 1, 1);
        }
    }
}

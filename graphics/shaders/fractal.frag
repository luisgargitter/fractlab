#version 330 core

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

float divBound(float alpha, float beta, vec2 prec) {
    alpha = abs(alpha);
    beta = abs(beta);
    float c = length(prec);
    float p2 = (beta+1)/(2*alpha);
    float q = -c/alpha;
    float R = p2 + sqrt(p2*p2 - q);

    return R;
}

vec3 colorFromHueSat(float hue, float sat) {
    vec3 cw[6] = vec3[](vec3(0, 0, 1), vec3(0, 1, 0), vec3(0, 1, 1), vec3(1, 0, 0), vec3(1, 0, 1), vec3(1, 1, 0));
    float h6 = hue*6.0;
    int i = int(h6);
    vec3 c1 = cw[i];
    vec3 c2 = cw[(i+1) % 6];
    float t = h6 - float(i);
    vec3 c = c1 * (1-t) + c2 * t;
    return sat * c/length(c);
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

    for(int i = 0; i < depth; i++) {
        y = PZn[0]*cMul(y, y) + PZn[1]*y + prec;

        if(length(y) < R) {
            j++;
        }
    }
    float hue = float(j)/float(depth);
    float sat = 1- pow(0.5 ,float(j));

    color = vec4(colorFromHueSat(hue, sat), 1.0);
}
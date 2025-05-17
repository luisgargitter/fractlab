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

const int depth = 64;
const int limitr = depth;
const int limitg = depth/8;
const int limitb = depth/64;

void main()
{
    vec2 x = coord;
    vec2 x2 = cMul(x, x);
    vec2 prec = PZ0[0]*x2 + PZ0[1]*x + c;
    float R = divBound(PZn[0], PZn[1], prec);

    vec2 y = x;

    int j = 0;
    int k = 0;
    int l = 0;

    for(int i = 0; i < depth; i++) {
        y = PZn[0]*cMul(y, y) + PZn[1]*y + prec;

        //y = cMul(y, y) + x; // mandelbrot set
        //z = cAdd(cMul(z, z), vec2(-0.5, 0.6)); // julia set

        if(length(y) < R) {
            if(i < limitr) j++;
            if(i < limitg) k++;
            if(i < limitb) l++;
        }
    }
    float r = float(j)/float(limitr);
    float g = float(k)/float(limitg);
    float b = float(l)/float(limitb);

    color = vec4(1-r, 1-g, 1-b, 1.0);
}
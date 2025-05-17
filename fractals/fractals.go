package fractals

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Polynomial = float32
type Fractal struct {
	C   complex64
	PZ0 [2]Polynomial
	PZn [2]Polynomial
}

func Mandelbrot() Fractal {
	return Fractal{
		C:   0.0 + 1i*0.0,
		PZ0: [2]Polynomial{0, 1},
		PZn: [2]Polynomial{1, 0},
	}
}

func Julia(c complex64) Fractal {
	return Fractal{
		C:   c,
		PZ0: [2]Polynomial{0, 0},
		PZn: [2]Polynomial{1, 0},
	}
}

type Animation struct {
	Src, Dest Fractal
	T         float32
}

func GetFractal(a Animation) Fractal {

	vsrc := mgl32.NewVecNFromData([]float32{real(a.Src.C), imag(a.Src.C), a.Src.PZ0[0], a.Src.PZ0[1], a.Src.PZn[0], a.Src.PZn[1]})
	vdest := mgl32.NewVecNFromData([]float32{real(a.Dest.C), imag(a.Dest.C), a.Dest.PZ0[0], a.Dest.PZ0[1], a.Dest.PZn[0], a.Dest.PZn[1]})
	vres := vsrc.Mul(nil, 1-a.T).Add(nil, vdest.Mul(nil, a.T)).Raw()
	return Fractal{
		C:   complex(vres[0], vres[1]),
		PZ0: [2]Polynomial{vres[2], vres[3]},
		PZn: [2]Polynomial{vres[4], vres[5]},
	}
}

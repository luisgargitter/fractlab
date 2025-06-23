package fractals

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Polynomial = float32
type Fractal struct {
	C   mgl32.Vec2
	PZ0 [2]Polynomial
	PZn [2]Polynomial
}

func Mandelbrot() Fractal {
	return Fractal{
		C:   mgl32.Vec2{0.0, 0.0},
		PZ0: [2]Polynomial{0, 1},
		PZn: [2]Polynomial{1, 0},
	}
}

func Julia(c complex64) Fractal {
	return Fractal{
		C:   mgl32.Vec2{real(c), imag(c)},
		PZ0: [2]Polynomial{0, 0},
		PZn: [2]Polynomial{1, 0},
	}
}

type Animation struct {
	Src, Dest Fractal
	Time      float32
}

func GetFractal(a Animation) Fractal {
	vsrc := mgl32.NewVecNFromData([]float32{a.Src.C[0], a.Src.C[1], a.Src.PZ0[0], a.Src.PZ0[1], a.Src.PZn[0], a.Src.PZn[1]})
	vdest := mgl32.NewVecNFromData([]float32{a.Dest.C[0], a.Dest.C[1], a.Dest.PZ0[0], a.Dest.PZ0[1], a.Dest.PZn[0], a.Dest.PZn[1]})
	vres := vsrc.Mul(nil, 1-a.Time).Add(nil, vdest.Mul(nil, a.Time)).Raw()
	return Fractal{
		C:   mgl32.Vec2{vres[0], vres[1]},
		PZ0: [2]Polynomial{vres[2], vres[3]},
		PZn: [2]Polynomial{vres[4], vres[5]},
	}
}

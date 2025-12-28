package objects

import "math"

type Vector2 struct {
	X, Y float64
}

type Vector3 struct {
	X, Y, Z float64
}

type Vector4 struct {
	X, Y, Z, W float64
}

var WorldRight = Vector3{1, 0, 0}
var WorldLeft = WorldRight.Invert()
var WorldUp = Vector3{0, -1, 0}
var WorldDown = WorldUp.Invert()
var WorldBackward = Vector3{0, 0, -1}
var WorldForward = WorldBackward.Invert()

func (v Vector2) Add(o Vector2) Vector2 {
	v.X += o.X
	v.Y += o.Y
	return v
}

func (v Vector2) Sub(o Vector2) Vector2 {
	v.X -= o.X
	v.Y -= o.Y
	return v
}

func (v Vector2) Mulf(f float64) Vector2 {
	v.X *= f
	v.Y *= f
	return v
}

func (v Vector2) Rotate(angle float64) Vector2 {
	tsin, tcos := math.Sincos(angle)
	x := v.X*tcos - v.Y*tsin
	y := v.X*tsin + v.Y*tcos
	v.X = x
	v.Y = y
	return v
}

func (v Vector2) Distance() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) DistanceSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) DistanceTo(o Vector2) float64 {
	return math.Sqrt(o.X-v.X)*(o.X-v.X) + (o.Y-v.Y)*(o.Y-v.Y)
}

func (v Vector2) DistanceSquaredTo(o Vector2) float64 {
	return (o.X-v.X)*(o.X-v.X) + (o.Y-v.Y)*(o.Y-v.Y)
}

func (v Vector2) Normalize() Vector2 {
	l := v.Distance()
	if l < 1e-8 || l == 1 {
		return v
	}
	v.X, v.Y = v.X/l, v.Y/l
	return v
}

func (v Vector2) Dot(o Vector2) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vector2) Cross(o Vector2) float64 {
	return v.X*o.Y - o.X*v.Y
}

func (v Vector3) Add(o Vector3) Vector3 {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
	return v
}

func (v Vector3) Sub(o Vector3) Vector3 {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
	return v
}

func (v Vector3) Invert() Vector3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() Vector3 {
	l := v.Magnitude()
	if l < 1e-8 || l == 1 {
		return v
	}
	v.X, v.Y, v.Z = v.X/l, v.Y/l, v.Z/l
	return v
}

func (v Vector3) Dot(o Vector3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector3) Cross(o Vector3) Vector3 {

	x := v.X
	y := v.Y
	z := v.Z

	v.X = y*o.Z - z*o.Y
	v.Y = z*o.X - x*o.Z
	v.Z = x*o.Y - y*o.X

	return v
}

func (v Vector3) RotateX(angle float64) Vector3 {
	a := angle
	y := v.Y
	z := v.Z

	v.Y = y*math.Cos(a) - z*math.Sin(a)
	v.Z = y*math.Sin(a) + z*math.Cos(a)
	return v
}

func (v Vector3) RotateY(angle float64) Vector3 {
	a := angle
	x := v.X
	z := v.Z

	v.X = x*math.Cos(a) + z*math.Sin(a)
	v.Z = -x*math.Sin(a) + z*math.Cos(a)
	return v
}

func (v Vector3) RotateZ(angle float64) Vector3 {
	a := angle
	x := v.X
	y := v.Y

	v.X = x*math.Cos(a) - y*math.Sin(a)
	v.Y = x*math.Sin(a) + y*math.Cos(a)
	return v
}

func (v Vector3) MulMat(m Matrix4) Vector3 {

	return Vector3{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0],
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1],
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2],
	}
}

func (v Vector3) MulMatToVector4(m Matrix4) Vector4 {

	return Vector4{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0],
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1],
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2],
		W: m[0][3]*v.X + m[1][3]*v.Y + m[2][3]*v.Z + m[3][3],
	}
}

func (v Vector3) Mulf(f float64) Vector3 {
	return Vector3{v.X * f, v.Y * f, v.Z * f}
}

func (v Vector4) MulMat(m Matrix4) Vector4 {

	return Vector4{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0]*v.W,
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1]*v.W,
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2]*v.W,
		W: m[0][3]*v.X + m[1][3]*v.Y + m[2][3]*v.Z + m[3][3]*v.W,
	}
}

func (v Vector3) Equals(o Vector3) bool {

	eps := 1e-4

	if math.Abs(v.X-o.X) > eps || math.Abs(v.Y-o.Y) > eps || math.Abs(v.Z-o.Z) > eps {
		return false
	}

	return true
}

func (v Vector4) ToVec3() Vector3 {
	return Vector3{v.X, v.Y, v.Z}
}

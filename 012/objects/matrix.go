package objects

import "math"

type Matrix4 [4][4]float64

func NewMatrix4() Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func NewEmptyMatrix4() Matrix4 {
	return Matrix4{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
}

func NewProjectionPerspective(zn, zf, left, right, top, bottom float64) Matrix4 {
	return Matrix4{
		{(2 * zn) / (right - left), 0, 0, 0},
		{0, (2 * zn) / (bottom - top), 0, 0},
		{0, 0, zf / (zf - zn), 1},
		{0, 0, -zn * zf / (zf - zn), 0},
	}
}

func NewProjectionOrthographic(zn, zf, left, right, top, bottom float64) Matrix4 {
	return Matrix4{
		{2 / (right - left), 0, 0, 0},
		{0, 2 / (top - bottom), 0, 0},
		{0, 0, 2 / (zf - zn), 0},
		{0, 0, 0, 1},
	}
}

func NewViewport(width, height float64) Matrix4 {
	return Matrix4{
		{width / 2, 0, 0, 0},
		{0, height / 2, 0, 0},
		{0, 0, 1, 0},
		{width / 2, height / 2, 0, 1},
	}
}

func NewMatrix4Translate(x, y, z float64) Matrix4 {
	mat := NewMatrix4()
	mat[3][0] = x
	mat[3][1] = y
	mat[3][2] = z
	return mat
}

func NewMatrix4Scale(x, y, z float64) Matrix4 {
	mat := NewMatrix4()
	mat[0][0] = x
	mat[1][1] = y
	mat[2][2] = z
	mat[3][3] = 1
	return mat
}

func NewMatrix4Rotate(x, y, z, angle float64) Matrix4 {

	if x == 0 && y == 0 && z == 0 {
		y = 1
	}

	mat := NewMatrix4()
	v := Vector3{X: x, Y: y, Z: z}.Normalize()
	s := math.Sin(angle)
	c := math.Cos(angle)
	m := 1 - c

	mat[0][0] = m*v.X*v.X + c
	mat[0][1] = m*v.X*v.Y + v.Z*s
	mat[0][2] = m*v.Z*v.X - v.Y*s

	mat[1][0] = m*v.X*v.Y - v.Z*s
	mat[1][1] = m*v.Y*v.Y + c
	mat[1][2] = m*v.Y*v.Z + v.X*s

	mat[2][0] = m*v.Z*v.X + v.Y*s
	mat[2][1] = m*v.Y*v.Z - v.X*s
	mat[2][2] = m*v.Z*v.Z + c

	return mat
}

func (m Matrix4) Mul(o Matrix4) Matrix4 {
	newM := NewMatrix4()

	newM[0][0] = m[0][0]*o[0][0] + m[0][1]*o[1][0] + m[0][2]*o[2][0] + m[0][3]*o[3][0]
	newM[1][0] = m[1][0]*o[0][0] + m[1][1]*o[1][0] + m[1][2]*o[2][0] + m[1][3]*o[3][0]
	newM[2][0] = m[2][0]*o[0][0] + m[2][1]*o[1][0] + m[2][2]*o[2][0] + m[2][3]*o[3][0]
	newM[3][0] = m[3][0]*o[0][0] + m[3][1]*o[1][0] + m[3][2]*o[2][0] + m[3][3]*o[3][0]

	newM[0][1] = m[0][0]*o[0][1] + m[0][1]*o[1][1] + m[0][2]*o[2][1] + m[0][3]*o[3][1]
	newM[1][1] = m[1][0]*o[0][1] + m[1][1]*o[1][1] + m[1][2]*o[2][1] + m[1][3]*o[3][1]
	newM[2][1] = m[2][0]*o[0][1] + m[2][1]*o[1][1] + m[2][2]*o[2][1] + m[2][3]*o[3][1]
	newM[3][1] = m[3][0]*o[0][1] + m[3][1]*o[1][1] + m[3][2]*o[2][1] + m[3][3]*o[3][1]

	newM[0][2] = m[0][0]*o[0][2] + m[0][1]*o[1][2] + m[0][2]*o[2][2] + m[0][3]*o[3][2]
	newM[1][2] = m[1][0]*o[0][2] + m[1][1]*o[1][2] + m[1][2]*o[2][2] + m[1][3]*o[3][2]
	newM[2][2] = m[2][0]*o[0][2] + m[2][1]*o[1][2] + m[2][2]*o[2][2] + m[2][3]*o[3][2]
	newM[3][2] = m[3][0]*o[0][2] + m[3][1]*o[1][2] + m[3][2]*o[2][2] + m[3][3]*o[3][2]

	newM[0][3] = m[0][0]*o[0][3] + m[0][1]*o[1][3] + m[0][2]*o[2][3] + m[0][3]*o[3][3]
	newM[1][3] = m[1][0]*o[0][3] + m[1][1]*o[1][3] + m[1][2]*o[2][3] + m[1][3]*o[3][3]
	newM[2][3] = m[2][0]*o[0][3] + m[2][1]*o[1][3] + m[2][2]*o[2][3] + m[2][3]*o[3][3]
	newM[3][3] = m[3][0]*o[0][3] + m[3][1]*o[1][3] + m[3][2]*o[2][3] + m[3][3]*o[3][3]

	return newM
}

func (matrix Matrix4) Inverted() Matrix4 {
	var A2323 = matrix[2][2]*matrix[3][3] - matrix[2][3]*matrix[3][2]
	var A1323 = matrix[2][1]*matrix[3][3] - matrix[2][3]*matrix[3][1]
	var A1223 = matrix[2][1]*matrix[3][2] - matrix[2][2]*matrix[3][1]
	var A0323 = matrix[2][0]*matrix[3][3] - matrix[2][3]*matrix[3][0]
	var A0223 = matrix[2][0]*matrix[3][2] - matrix[2][2]*matrix[3][0]
	var A0123 = matrix[2][0]*matrix[3][1] - matrix[2][1]*matrix[3][0]
	var A2313 = matrix[1][2]*matrix[3][3] - matrix[1][3]*matrix[3][2]
	var A1313 = matrix[1][1]*matrix[3][3] - matrix[1][3]*matrix[3][1]
	var A1213 = matrix[1][1]*matrix[3][2] - matrix[1][2]*matrix[3][1]
	var A2312 = matrix[1][2]*matrix[2][3] - matrix[1][3]*matrix[2][2]
	var A1312 = matrix[1][1]*matrix[2][3] - matrix[1][3]*matrix[2][1]
	var A1212 = matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]
	var A0313 = matrix[1][0]*matrix[3][3] - matrix[1][3]*matrix[3][0]
	var A0213 = matrix[1][0]*matrix[3][2] - matrix[1][2]*matrix[3][0]
	var A0312 = matrix[1][0]*matrix[2][3] - matrix[1][3]*matrix[2][0]
	var A0212 = matrix[1][0]*matrix[2][2] - matrix[1][2]*matrix[2][0]
	var A0113 = matrix[1][0]*matrix[3][1] - matrix[1][1]*matrix[3][0]
	var A0112 = matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]

	var det = matrix[0][0]*(matrix[1][1]*A2323-matrix[1][2]*A1323+matrix[1][3]*A1223) -
		matrix[0][1]*(matrix[1][0]*A2323-matrix[1][2]*A0323+matrix[1][3]*A0223) +
		matrix[0][2]*(matrix[1][0]*A1323-matrix[1][1]*A0323+matrix[1][3]*A0123) -
		matrix[0][3]*(matrix[1][0]*A1223-matrix[1][1]*A0223+matrix[1][2]*A0123)

	det = 1 / det

	m := NewMatrix4()

	m[0][0] = det * (matrix[1][1]*A2323 - matrix[1][2]*A1323 + matrix[1][3]*A1223)
	m[0][1] = det * -(matrix[0][1]*A2323 - matrix[0][2]*A1323 + matrix[0][3]*A1223)
	m[0][2] = det * (matrix[0][1]*A2313 - matrix[0][2]*A1313 + matrix[0][3]*A1213)
	m[0][3] = det * -(matrix[0][1]*A2312 - matrix[0][2]*A1312 + matrix[0][3]*A1212)
	m[1][0] = det * -(matrix[1][0]*A2323 - matrix[1][2]*A0323 + matrix[1][3]*A0223)
	m[1][1] = det * (matrix[0][0]*A2323 - matrix[0][2]*A0323 + matrix[0][3]*A0223)
	m[1][2] = det * -(matrix[0][0]*A2313 - matrix[0][2]*A0313 + matrix[0][3]*A0213)
	m[1][3] = det * (matrix[0][0]*A2312 - matrix[0][2]*A0312 + matrix[0][3]*A0212)
	m[2][0] = det * (matrix[1][0]*A1323 - matrix[1][1]*A0323 + matrix[1][3]*A0123)
	m[2][1] = det * -(matrix[0][0]*A1323 - matrix[0][1]*A0323 + matrix[0][3]*A0123)
	m[2][2] = det * (matrix[0][0]*A1313 - matrix[0][1]*A0313 + matrix[0][3]*A0113)
	m[2][3] = det * -(matrix[0][0]*A1312 - matrix[0][1]*A0312 + matrix[0][3]*A0112)
	m[3][0] = det * -(matrix[1][0]*A1223 - matrix[1][1]*A0223 + matrix[1][2]*A0123)
	m[3][1] = det * (matrix[0][0]*A1223 - matrix[0][1]*A0223 + matrix[0][2]*A0123)
	m[3][2] = det * -(matrix[0][0]*A1213 - matrix[0][1]*A0213 + matrix[0][2]*A0113)
	m[3][3] = det * (matrix[0][0]*A1212 - matrix[0][1]*A0212 + matrix[0][2]*A0112)

	return m

}

func NewMatrix4LookAt(from, to, up Vector3) Matrix4 {

	if from.Equals(to) {
		return NewMatrix4()
	}
	z := to.Sub(from).Normalize().Invert()

	up = up.Normalize()

	if z.Equals(up) || z.Equals(up.Invert()) {
		if !up.Equals(WorldRight) {
			up = WorldRight
		} else {
			up = WorldBackward
		}
	}

	x := z.Cross(up).Normalize()
	y := z.Cross(x)
	xt := -x.Dot(from)
	yt := -y.Dot(from)
	zt := -z.Dot(from)
	return Matrix4{
		{x.X, y.X, z.X, 0},
		{x.Y, -y.Y, z.Y, 0},
		{x.Z, y.Z, z.Z, 0},
		{xt, yt, zt, 1},
	}
}

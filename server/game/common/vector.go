package common

import "math"

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func VectorSum(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func VectorDotProduct(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X * v2.X, Y: v1.Y * v2.Y}
}

func (v Vector2) Magnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) SqrtMagnitude() float64 {
	return math.Sqrt(v.Magnitude())
}

func (v Vector2) ScalarProduct(scalar float64) Vector2 {
	v2 := Vector2{v.X * scalar, v.Y * scalar}
	return v2
}

func (v Vector2) ClampVector(minX, maxX, minY, maxY float64) Vector2 {
	v2 := Vector2{v.X, v.Y}
	v2.X = Clamp(v2.X, minX, maxX)
	v2.Y = Clamp(v2.Y, minY, maxY)

	return v2
}

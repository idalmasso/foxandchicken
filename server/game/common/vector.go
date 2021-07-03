package common

type Vector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func VectorSum(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func VectorDotProduct(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X * v2.X, Y: v1.Y * v2.Y}
}

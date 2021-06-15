package game

type Vector2 struct {
	X, Y float32
}

func VectorSum(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func VectorDotProduct(v1, v2 Vector2) Vector2 {
	return Vector2{X: v1.X * v2.X, Y: v1.Y * v2.Y}
}

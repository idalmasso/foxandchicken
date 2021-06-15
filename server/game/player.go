package game

type PlayerGameData struct {
	CharacterType int
	Username      string `json:"username"`
	Position      Vector2
	Rotation      float32
	Velocity      Vector2
	SizeRadius    float32
}

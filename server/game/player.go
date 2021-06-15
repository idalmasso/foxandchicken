package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Username string          `json:"username"`
	GameData *PlayerGameData `json:"data"`
	mutex    sync.Mutex
	Conn     *websocket.Conn
}

type PlayerGameData struct {
	CharacterType int
	Position      Vector2
	Rotation      float32
	Velocity      Vector2
	SizeRadius    float32
}

func (p *Player) UpdateWebSocket(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Conn = conn
}

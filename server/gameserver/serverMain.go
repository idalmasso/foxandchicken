package gameserver

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/foxandchicken/server/game"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//user is the login type from json
type user struct {
	Username string `json:"username"`
}

//GameServer is a game server with rooms and other things...
type GameServer struct {
	Instance *game.GameInstance
}

func (gameServer GameServer) ManageRequest(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := NewPlayer(gameServer.Instance)

	if err != nil {
		ws.Close()
		return
	}
	p.UpdateWebSocket(ws)
}

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

type user struct {
	Username string `json:"username"`
}
type userResponse struct {
	Id string `json:"id"`
}
type GameServer struct {
	Instance *game.GameInstance
}

func (gameServer GameServer) ManageRequest(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var p Player
	p.GameInstance = gameServer.Instance

	if err != nil {
		ws.Close()
		return
	}
	p.UpdateWebSocket(ws)
}

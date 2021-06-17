package gameserver

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/foxandchicken/server/game"
)

type Player struct {
	GameData     *game.PlayerGameData `json:"data"`
	GameInstance *game.GameInstance
	mutex        sync.Mutex
	Conn         *websocket.Conn
}

func (p *Player) UpdateWebSocket(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Conn = conn
	go p.PlayerCycle()
}

func (p *Player) PlayerCycle() {
	p.ReadUsername()

	type message struct {
		Message string `json:"message"`
	}
	var mex message
	for {
		if err := p.Conn.ReadJSON(&mex); err != nil {

			fmt.Println("ERROR "+p.GameData.Username, "Cannot decode the chat message", err.Error())
			p.Conn.Close()
			p.GameInstance.RemovePlayer(p.GameData.Username)
			return
		}
		fmt.Println("Received message " + mex.Message + " from user " + p.GameData.Username)
	}

}

//ReadUsername block the user until an ok username is inserted
func (p *Player) ReadUsername() {
	type usernameMessage struct {
		Username string `json:"username"`
	}
	var u usernameMessage
	ok := false
	for !ok {
		err := p.Conn.ReadJSON(&u)
		if err != nil {
			p.Conn.WriteJSON(struct{ message string }{message: "Error: " + err.Error()})
			p.Conn.Close()
			return
		}
		if u.Username != "" {
			p.GameData, err = p.GameInstance.AddPlayer(u.Username)
			if err == nil {
				ok = true
				p.Conn.WriteJSON(struct{ message string }{message: "OK"})
			} else {
				p.Conn.WriteJSON(struct{ message string }{message: "Error: " + err.Error()})
			}
		} else {
			p.Conn.WriteJSON(struct{ message string }{message: "Error: not empty"})
		}
	}
}

package gameserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/foxandchicken/server/game"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type actionMessageTypes string

const (
	ActionMessageCreateRoom actionMessageTypes = "CREATEROOM"
	ActionMessageLeaveRoom  actionMessageTypes = "LEAVEROOM"
	ActionMessageJoinRoom   actionMessageTypes = "JOINROOM"
)

type usernameMessage struct {
	Username string `json:"username"`
}

type message struct {
	Action  actionMessageTypes `json:"action"`
	Message string             `json:"message"`
}

type singleStringReturnMessage struct {
	Message string `json:"message"`
}

type Player struct {
	GameData     *game.PlayerGameData `json:"data"`
	GameInstance *game.GameInstance
	RoomChannel  chan<- messaging.RoomMessageValue
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
	var mex message
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err := p.Conn.ReadJSON(&mex); err != nil {
			fmt.Println("ERROR "+p.GameData.Username, "cannot decode the message", err.Error())
			p.Conn.Close()
			p.GameInstance.RemovePlayer(p.GameData.Username)
			return
		}
		switch mex.Action {
		case ActionMessageCreateRoom:
			if err := p.tryCreateRoom(mex.Message); err != nil {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
			}
		case ActionMessageLeaveRoom:
			if err := p.tryLeaveRoom(); err != nil {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
			}
		default:
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
		}

		fmt.Println("Received message " + mex.Message + " from user " + p.GameData.Username)
	}

}

//ReadUsername block the user until an ok username is inserted
func (p *Player) ReadUsername() {
	var u usernameMessage
	ok := false
	for !ok {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		err := p.Conn.ReadJSON(&u)
		if err != nil {
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
		} else if u.Username != "" {
			p.GameData, err = p.GameInstance.AddPlayer(u.Username)
			if err == nil {
				ok = true
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "requested message: username:<'username'> error: " + err.Error()})
			}
		} else {
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "requested message: username:<'username'> error:  empty username"})
		}

	}
}

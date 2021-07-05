package gameserver

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/foxandchicken/server/game"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type usernameMessage struct {
	Username string `json:"username"`
}

type Player struct {
	username          string
	GameData          *game.PlayerGameData `json:"data"`
	GameInstance      *game.GameInstance
	RoomChannel       chan<- messaging.RoomMessageValue
	RoomChannelOutput <-chan messaging.RoomMessageValue
	IsInRoom          bool
	mutex             sync.Mutex
	Conn              *websocket.Conn
	EndGameChannel    chan bool
	EndPlayer         chan bool
}

func (p *Player) UpdateWebSocket(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Conn = conn
	go p.PlayerCycle()

}

func (p *Player) PlayerCycle() {
	p.ReadUsername()
	go p.PlayerBroadcastListener()
	var mex message
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err := p.Conn.ReadJSON(&mex); err != nil {
			log.Println("ERROR "+p.username, "cannot decode the message", err.Error())
			p.Conn.Close()
			p.GameInstance.RemovePlayer(p.username)
			p.EndPlayer <- true
			return
		}
		fmt.Println("Received message " + mex.Message + " from user " + p.username)
		switch mex.Action {
		case ActionMessageCreateRoom:
			if err := p.tryCreateRoom(mex.Message); err != nil {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
				p.mutex.Unlock()
			} else {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
				p.mutex.Unlock()
				if err := p.PlayerRoomInputCycle(); err != nil {
					p.EndPlayer <- true
					return
				}

			}
		case ActionMessageJoinRoom:
			if err := p.tryJoinRoom(mex.Message); err != nil {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
				p.mutex.Unlock()
			} else {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
				p.mutex.Unlock()
				if err := p.PlayerRoomInputCycle(); err != nil {
					p.EndPlayer <- true
					return
				}
			}
		default:
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			p.mutex.Unlock()

		}

	}
}

func (p *Player) tryCreateRoom(roomName string) error {
	var m messaging.CommMessageCreateRoom
	m.Player = p.username
	m.Name = roomName
	v, err := p.sendAndReturnError(&m, messaging.MessageResponseCreateRoom)
	if err != nil {
		return err
	}
	ret := v.(*messaging.CommMessageResponseCreateRoom)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RoomChannel = ret.RoomChannel
	p.IsInRoom = true
	return nil
}
func (p *Player) tryJoinRoom(roomName string) error {
	var m messaging.CommRoomMessageJoinPlayer
	m.Player = p.username
	m.Name = roomName
	v, err := p.sendAndReturnError(&m, messaging.MessageResponseJoinRoom)
	if err != nil {
		return err
	}
	ret := v.(*messaging.CommMessageResponseJoinRoom)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RoomChannel = ret.RoomChannel
	p.IsInRoom = true
	return nil
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
			err = p.GameInstance.AddPlayer(u.Username)
			if err == nil {
				ok = true
				p.mutex.Lock()
				p.username = u.Username

				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
				p.mutex.Unlock()
			} else {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "requested message: username:<'username'> error: " + err.Error()})
				p.mutex.Unlock()
			}
		} else {
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "requested message: username:<'username'> error:  empty username"})
			p.mutex.Unlock()
		}

	}
}
func (p *Player) sendAndReturnError(m messaging.InstanceMessageValue, acceptedType messaging.MessageType) (messaging.InstanceMessageValue, error) {
	p.GameInstance.InputChannel <- m
	v := <-p.GameInstance.PlayerDataChannels[p.username]
	if v.GetMessageType() != acceptedType {
		return nil, fmt.Errorf("wrong message type in return")
	}
	if v.ErrorMessage() != "" {
		return nil, fmt.Errorf(v.ErrorMessage())
	}

	return v, nil
}

func NewPlayer(instance *game.GameInstance) *Player {
	var p Player
	p.GameInstance = instance
	p.IsInRoom = false
	p.EndGameChannel = make(chan bool)

	return &p
}
func (p *Player) PlayerBroadcastListener() {
	for {
		select {
		case <-p.EndPlayer:
			log.Println("Player Broadcast exit" + p.username)
			return
		case m := <-p.GameInstance.PlayerDataChannelsBroadcasts[p.username]:
			switch m.GetMessageType() {
			default:
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "got message broadcast" + m.ErrorMessage()})
				p.mutex.Unlock()
			}
		}
	}
}

package gameserver

import (
	"encoding/json"
	"errors"
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
	IsClosing         bool
}

//UpdateWebSocket updates the websocket connection in the player
func (p *Player) UpdateWebSocket(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Conn = conn
	go p.PlayerCycle()

}

//PlayerCycle is the cycle of a player when not in the room
func (p *Player) PlayerCycle() {
	if err := p.ReadUsername(); err != nil {
		return
	}
	go p.PlayerBroadcastListener()
	var mex message
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err := p.Conn.ReadJSON(&mex); err != nil {
			var jErr *json.SyntaxError
			if errors.As(err, &jErr) {
				log.Println("ERROR "+p.username, "cannot decode the message", err.Error())
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: TIMEOUT"})
				p.Close()
				p.GameInstance.RemovePlayer(p.username)
				return
			}
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
		case ActionListRooms:
			rooms := p.GameInstance.GetRooms()
			p.mutex.Lock()
			p.Conn.WriteJSON(rooms)
			p.mutex.Unlock()
		default:
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			p.mutex.Unlock()

		}

	}
}

//tryCreateRoom tries to create a named room in the server. Automatically joins it
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
	p.RoomChannelOutput = ret.RoomResponseChannel
	p.IsInRoom = true
	return nil
}

//tryJoinRoom tries to join a named room in the server
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
	p.RoomChannelOutput = ret.RoomResponseChannel
	p.IsInRoom = true
	return nil
}

//ReadUsername block the user until an ok username is inserted
func (p *Player) ReadUsername() error {
	var u usernameMessage
	ok := false
	for !ok {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		err := p.Conn.ReadJSON(&u)
		if err != nil {
			var jErr *json.SyntaxError
			if errors.As(err, &jErr) {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: TIMEOUT"})
				p.Conn.Close()
				return fmt.Errorf("TIMEOUT")
			}

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
	return nil
}

//sendAndReturnError send a message to the instance and test its return value
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

//NewPlayer returns a new Player using a gameInstance instance
func NewPlayer(instance *game.GameInstance) *Player {
	var p Player
	p.GameInstance = instance
	p.IsInRoom = false
	p.EndGameChannel = make(chan bool)
	p.IsClosing = false
	return &p
}

//PlayerBroadcastListener is the listener for the broadcast messages from the gameInstance
func (p *Player) PlayerBroadcastListener() {
	for {
		select {
		case <-p.EndPlayer:
			log.Println("Player Broadcast exit" + p.username)
			close(p.EndPlayer)
			return
		case m := <-p.GameInstance.PlayerDataChannelsBroadcasts[p.username]:
			p.mutex.Lock()
			if !p.IsClosing {
				switch m.GetMessageType() {
				default:
					p.mutex.Lock()
					p.Conn.WriteJSON(singleStringReturnMessage{Message: "got message broadcast" + m.ErrorMessage()})
					p.mutex.Unlock()
				}
			}
			p.mutex.Unlock()
		}
	}
}

//Close close the player handles
func (p *Player) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.IsClosing=true
	if p.IsInRoom {
		p.EndGameChannel <- true
		
	}
	
	p.EndPlayer <- true
	p.Conn.Close()

}

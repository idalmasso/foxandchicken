package gameserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

//sendAndReturnErrorRoom send a message to the room and wait and test a return value type
func (p *Player) sendAndReturnErrorRoom(m messaging.RoomMessageValue, acceptedType messaging.MessageType) (messaging.RoomMessageValue, error) {
	p.RoomChannel <- m
	v := <-p.GameInstance.PlayerDataChannels[p.username]
	if v.GetMessageType() != acceptedType {
		return nil, fmt.Errorf("wrong message type in return")
	}
	if v.ErrorMessage() != "" {
		return nil, fmt.Errorf(v.ErrorMessage())
	}
	return nil, nil
}

//PlayerRoomInputCycle is the room cycle for a player. (after join a room this is the cycle for input)
func (p *Player) PlayerRoomInputCycle() error {
	var mex genericMessage
	go p.PlayerRoomGameCycle()
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err := p.Conn.ReadJSON(&mex); err != nil {
			var jErr *json.SyntaxError
			if errors.As(err, &jErr) {
				log.Println("ERROR "+p.username, "cannot decode the message", err.Error())
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
			} else {
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: TIMEOUT"})
				p.Conn.Close()
				p.GameInstance.RemovePlayer(p.username)
				return err
			}
		}
		if !p.IsInRoom {
			return nil
		}

		switch mex.GetAction() {
		case ActionMessageLeaveRoom:
			if err := p.tryLeaveRoom(); err != nil {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
				p.EndGameChannel <- true
				p.mutex.Unlock()
				return nil
			} else {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
				p.EndGameChannel <- true
				p.mutex.Unlock()
				return nil
			}
		case ActionMessageMovement:
			jsonString, _ := json.Marshal(mex.Message)
			m := movementStruct{}
			err := json.Unmarshal(jsonString, &m)
			//m, ok := mex.Message.(movementStruct)

			if err != nil {
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "message not recognized"})
				p.mutex.Unlock()
			} else {
				p.mutex.Lock()
				position := common.Vector2{X: m.PositionX, Y: m.PositionY}
				velocity := common.Vector2{X: m.VelocityX, Y: m.VelocityY}
				acceleration := common.Vector2{X: m.AccelerationX, Y: m.AccelerationY}
				sMex := messaging.CommRoomMessageMovePlayer{Player: p.username, Position: position, Velocity: velocity, Rotation: m.Rotation, Acceleration: acceleration}
				p.RoomChannel <- &sMex
				p.mutex.Unlock()
			}
		default:
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			p.mutex.Unlock()
		}

	}
}

//PlayerRoomGameCycle is the cycle for a single player that gets the messages from the server and write the message to the user
func (p *Player) PlayerRoomGameCycle() {
	for {
		select {
		case <-p.EndGameChannel:
			return
		case v := <-p.RoomChannelOutput:
			if !p.IsClosing {
				switch v.GetMessageType() {
				case messaging.RoomMessageTypePlayersMovement:
					p.mutex.Lock()
					moves := v.(*messaging.CommRoomMessagePlayersMovement)
					//log.Println("received message move>", move.Player)
					p.Conn.WriteJSON(moves)
					p.mutex.Unlock()
				case messaging.RoomMessageTypeLeftPlayer:
					p.mutex.Lock()
					mex := v.(*messaging.CommRoomMessageLeftPlayer)
					//log.Println("received message move>", move.Player)
					p.Conn.WriteJSON(message{Action: "LEAVEROOM", Message: mex.Player})
					p.mutex.Unlock()
				}
			}
		}
	}
}

func (p *Player) tryLeaveRoom() error {
	var m messaging.CommRoomMessageLeftPlayer
	m.Player = p.username

	_, err := p.sendAndReturnError(&m, messaging.MessageResponse)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.IsInRoom = false
	return err
}

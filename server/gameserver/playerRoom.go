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
				log.Println(p.username, "Game server Lock no decode")
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
				log.Println(p.username, "Game server unlock no decode")
				p.mutex.Unlock()
			} else {
				log.Println(p.username, "Timeout")
				log.Println(p.username, "Game server Lock timeout")
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: TIMEOUT"})
				log.Println(p.username, "Game server unlock timeout")
				p.mutex.Unlock()
				
				log.Println(p.username, "End of player room cycle")
				return err
			}
		}
		if !p.IsInRoom {
			return nil
		}

		switch mex.GetAction() {
		case ActionMessageLeaveRoom:
			if err := p.tryLeaveRoom(); err != nil {
				log.Println(p.username, "Game server Lock")
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: err.Error()})
				p.EndGameChannel <- true
				log.Println(p.username, "Game server Unlock")
				p.mutex.Unlock()
				log.Println(p.username, "Leave room - End of player room cycle WITH ERROR", err.Error())
				return err
			} else {
				log.Println(p.username, "Game server Lock ok leave room")
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "OK"})
				p.EndGameChannel <- true
				log.Println(p.username, "Game server Unlock ok leave room")
				p.mutex.Unlock()
				log.Println(p.username, "Leave room - End of player room cycle no error")
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
				position := common.Vector2{X: m.PositionX, Y: m.PositionY}
				velocity := common.Vector2{X: m.VelocityX, Y: m.VelocityY}
				acceleration := common.Vector2{X: m.AccelerationX, Y: m.AccelerationY}
				sMex := messaging.CommRoomMessageMovePlayer{Player: p.username, Position: position, Velocity: velocity, Rotation: m.Rotation, Acceleration: acceleration}
				log.Println(p.username, "Send message movement")
				p.RoomChannel <- &sMex
			}
		default:
			log.Println(p.username, "Game server Lock no rec")
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			log.Println(p.username, "Game server Unlock no rec")
			p.mutex.Unlock()
		}

	}
}

//PlayerRoomGameCycle is the cycle for a single player that gets the messages from the server and write the message to the user
func (p *Player) PlayerRoomGameCycle() {
	for {
		select {
		case <-p.EndGameChannel:
			close(p.EndGameChannel)
			log.Println(p.username, "Player end of PlayerRoomGameCycle")
			return
		case v := <-p.RoomChannelOutput:
			log.Println(p.username, "Game server Lock")
			p.mutex.Lock()
			if !p.IsClosing {
				switch v.GetMessageType() {
				case messaging.RoomMessageTypePlayersMovement:
					moves := v.(*messaging.CommRoomMessagePlayersMovement)
					//log.Println("received message move>", move.Player)
					p.Conn.WriteJSON(moves)
				case messaging.RoomMessageTypeLeftPlayer:
					mex := v.(*messaging.CommRoomMessageLeftPlayer)
					//log.Println("received message move>", move.Player)
					p.Conn.WriteJSON(message{Action: "LEAVEROOM", Message: mex.Player})
				}
			}
			log.Println(p.username, "Game server UnLock")
			p.mutex.Unlock()

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

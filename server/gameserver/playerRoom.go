package gameserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
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
				if glog.V(1) {
					glog.Warningln("Player.PlayerRoomInputCycle - ERROR "+p.username, "cannot decode the message", err.Error())
				}
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex lock no decode", p.username)
				}

				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: " + err.Error()})
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex unlock no decode", p.username)
				}
				p.mutex.Unlock()
			} else {
				if glog.V(1) {
					glog.Infoln("Player.PlayerRoomInputCycle - User", p.username, "Timeout")
				}
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex lock timeout", p.username)
				}
				p.mutex.Lock()
				p.Conn.WriteJSON(singleStringReturnMessage{Message: "error: TIMEOUT"})
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex unlock timeout", p.username)
				}
				p.mutex.Unlock()

				if glog.V(2) {
					glog.Infoln("Player.PlayerRoomInputCycle - end", p.username)
				}
				return err
			}
		}
		if !p.IsInRoom {
			return nil
		}

		switch mex.GetAction() {
		case ActionMessageLeaveRoom:
			if err := p.tryLeaveRoom(); err != nil {
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex lock ActionMessageLeaveRoom with error", p.username)
				}
				p.mutex.Lock()
				p.IsInRoom = false
				p.Conn.WriteJSON(genericMessage{Action: ActionMessageLeaveRoomResponse, Message: usernameErrorMessage{Username: p.username, Error: err.Error()}})
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex unlock ActionMessageLeaveRoom with error", p.username)
				}
				p.mutex.Unlock()
				if glog.V(1) {
					glog.Warningln("Player.PlayerRoomInputCycle - Leave room WITH error", p.username, err.Error())
				}
				return err
			} else {
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex lock ActionMessageLeaveRoom no error", p.username)
				}
				p.mutex.Lock()
				p.Conn.WriteJSON(genericMessage{Action: ActionMessageLeaveRoomResponse, Message: usernameErrorMessage{Username: p.username, Error: ""}})
				p.IsInRoom = false
				p.EndGameChannel <- true
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - p mutex unlock ActionMessageLeaveRoom no error", p.username)
				}
				p.mutex.Unlock()
				if glog.V(2) {
					glog.Infoln("Player.PlayerRoomInputCycle - Leave room no error", p.username)
				}
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
				sMex := messaging.CommRoomMessageMovePlayer{Player: p.username, ActionPressed: m.ActionPressed, Position: position, Velocity: velocity, Rotation: m.Rotation, Acceleration: acceleration}
				if glog.V(4) {
					glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - ActionMessageMovement", p.username, sMex)
				}
				p.RoomChannel <- &sMex
			}
		default:
			if glog.V(4) {
				glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - lock no recognited", p.username)
			}
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			if glog.V(4) {
				glog.Infoln("DEBUG - Player.PlayerRoomInputCycle - unlock no recognited", p.username)
			}
			p.mutex.Unlock()
		}

	}
}

//PlayerRoomGameCycle is the cycle for a single player that gets the messages from the server and write the message to the user
func (p *Player) PlayerRoomGameCycle() {
	if glog.V(2) {
		glog.Infoln("Player.PlayerRoomGameCycle - start", p.username)
	}
	for {
		select {
		case <-p.EndGameChannel:
			if glog.V(2) {
				glog.Infoln("Player.PlayerRoomGameCycle - end", p.username)
			}
			return
		case v := <-p.RoomChannelOutput:
			if glog.V(4) {
				glog.Infoln("DEBUG - Player.PlayerRoomGameCycle - RoomChannelOutput lock", p.username)
			}
			p.mutex.Lock()
			if !p.IsClosing {
				switch v.GetMessageType() {
				case messaging.RoomMessageTypePlayersStatuses:
					moves := v.(*messaging.CommRoomMessagePlayersStatuses)
					//log.Println("received message move>", move.Player)
					if data, err := json.Marshal(moves); err == nil {
						p.Conn.WriteJSON(message{Action: ActionMessageMovesRoom, Message: string(data)})
					} else {
						if glog.V(1) {
							glog.Warningln("Player.PlayerRoomGameCycle - RoomChannelOutput error json", err.Error())
						}
					}

				case messaging.RoomMessageTypeLeftPlayer:
					mex := v.(*messaging.CommRoomMessageLeftPlayer)
					//log.Println("received message move>", move.Player)
					p.Conn.WriteJSON(message{Action: ActionMessageLeaveRoom, Message: mex.Player})
				}
			}
			if glog.V(4) {
				glog.Infoln("DEBUG - Player.PlayerRoomGameCycle - RoomChannelOutput unlock", p.username)
			}
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

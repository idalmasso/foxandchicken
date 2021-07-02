package gameserver

import (
	"fmt"
	"log"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type movemementMessage struct {
	Action    actionMessageTypes `json:"action"`
	PositionX float32            `json:"position_x"`
	PositionY float32            `json:"position_y"`
	VelocityX float32            `json:"velocity_x"`
	VelocityY float32            `json:"velocity_y"`
	Rotation  float32            `json:"rotation"`
}

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

func (p *Player) PlayerRoomInputCycle() error {
	var mex message
	go p.PlayerRoomGameCycle()
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err := p.Conn.ReadJSON(&mex); err != nil {
			log.Println("ERROR "+p.username, "cannot decode the message", err.Error())

			p.Conn.Close()
			p.GameInstance.RemovePlayer(p.username)
			return err
		}
		if !p.IsInRoom {
			return nil
		}
		switch mex.Action {
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
		default:
			p.mutex.Lock()
			p.Conn.WriteJSON(singleStringReturnMessage{Message: "action not recognized"})
			p.mutex.Unlock()
		}
	}
}

func (p *Player) PlayerRoomGameCycle() {
	for {
		select {
		case <-p.EndGameChannel:
			return
		case v := <-p.RoomChannelOutput:
			if v.GetMessageType() == messaging.RoomMessageTypeMovePlayer {
				p.mutex.Lock()
				move := v.(*messaging.CommRoomMessageMovePlayer)
				log.Println("received message move>", move.Player)
				p.RoomChannel <- &messaging.CommRoomMessageMovePlayer{Player: p.username, Position: common.Vector2{X: move.Position.X, Y: move.Position.Y}}
				p.mutex.Unlock()

			}

		}
	}
}

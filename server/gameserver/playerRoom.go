package gameserver

import (
	"fmt"
	"log"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)



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
	var mex actionMessage
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
			m:=mex.(*movemementMessage)
			p.mutex.Lock()
			position := common.Vector2{X: m.PositionX, Y: m.PositionY}
			velocity := common.Vector2{X: m.VelocityX, Y: m.VelocityY}
			
			sMex:=messaging.CommRoomMessageMovePlayer{Player: p.username,Position: position, Velocity: velocity , Rotation: m.Rotation}
			p.RoomChannel<-&sMex
			p.mutex.Unlock()
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
			close(p.RoomChannel)
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

func (p *Player) tryLeaveRoom() error {
	var m messaging.CommRoomMessageLeftPlayer
	m.Player = p.username
	_, err := p.sendAndReturnError(&m, messaging.MessageResponse)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.IsInRoom = false
	return err
}

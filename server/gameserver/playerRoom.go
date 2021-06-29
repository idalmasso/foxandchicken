package gameserver

import (
	"fmt"

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

func (p *Player) tryCreateRoom(roomName string) error {
	var m messaging.CommMessageCreateRoom
	m.Player = p.GameData.Username
	m.Name = roomName
	v, err := p.sendAndReturnError(&m, messaging.MessageResponseCreateRoom)
	if err != nil {
		return err
	}
	ret := v.(*messaging.CommMessageResponseCreateRoom)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RoomChannel = ret.RoomChannel
	return nil
}
func (p *Player) tryJoinRoom(roomName string) error {
	var m messaging.CommRoomMessageJoinPlayer
	m.Player = p.GameData.Username
	m.Name = roomName
	v, err := p.sendAndReturnError(&m, messaging.MessageResponseCreateRoom)
	if err != nil {
		return err
	}
	ret := v.(*messaging.CommMessageResponseCreateRoom)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RoomChannel = ret.RoomChannel
	return nil
}
func (p *Player) tryLeaveRoom() error {
	var m messaging.CommRoomMessageLeftPlayer
	m.Player = p.GameData.Username
	_, err := p.sendAndReturnErrorRoom(&m, messaging.RoomMessageTypeResponseMessage)
	return err
}

func (p *Player) sendAndReturnError(m messaging.MessageValue, acceptedType messaging.MessageType) (messaging.MessageValue, error) {
	p.GameInstance.InputChannel <- m
	v := <-p.GameInstance.PlayerDataChannels[p.GameData.Username]
	if v.GetMessageType() != acceptedType {
		return nil, fmt.Errorf("wrong message type in return")
	}
	if v.ErrorMessage() != "" {
		return nil, fmt.Errorf(v.ErrorMessage())
	}

	return v, nil
}
func (p *Player) sendAndReturnErrorRoom(m messaging.MessageValue, acceptedType messaging.MessageType) (messaging.MessageValue, error) {
	p.RoomChannel <- m
	v := <-p.GameInstance.PlayerDataChannels[p.GameData.Username]
	if v.GetMessageType() != acceptedType {
		return nil, fmt.Errorf("wrong message type in return")
	}
	if v.ErrorMessage() != "" {
		return nil, fmt.Errorf(v.ErrorMessage())
	}
	return nil, nil
}

func (p *Player) PlayerRoomGameCycle() {

}

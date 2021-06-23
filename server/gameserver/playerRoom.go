package gameserver

import (
	"fmt"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

func (p *Player) tryCreateRoom(roomName string) error {
	var m messaging.CommMessageCreateRoom
	m.Player=p.GameData.Username
	m.Name=roomName
	return p.sendAndReturnError(&m)
}

func (p *Player) tryDeleteRoom(roomName string) error {
	var m messaging.CommMessageDeleteRoom
	m.Player=p.GameData.Username
	m.Name=roomName
	return p.sendAndReturnError(&m)
}

func (p *Player) sendAndReturnError(m messaging.MessageValue) error {
	p.GameInstance.InputChannel<-m
	v:=<-p.GameInstance.PlayerDataChannels[p.GameData.Username]
	if v.GetMessageType()!=messaging.MessageOkOrError{
		return fmt.Errorf("wrong message type in return from create room")
	}
	ret:=v.(*messaging.CommMessageOkOrError)
	if ret.Message!=""{
		return fmt.Errorf(ret.Message)
	}
	return nil
}

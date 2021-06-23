package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameInstance struct {
	Rooms              map[string]GameRoom
	Players            map[string]*PlayerGameData
	mutex              sync.Mutex
	waitingRoom        *WaitingRoom
	PlayerDataChannels map[string]chan messaging.MessageValue
	InputChannel       chan messaging.MessageValue
}

func (instance *GameInstance) AddPlayer(username string) (*PlayerGameData, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if _, ok := instance.Players[username]; ok {
		return nil, fmt.Errorf("already exists")
	}
	var p PlayerGameData
	p.Username = username
	instance.Players[username] = &p
	instance.PlayerDataChannels[username] = make(chan messaging.MessageValue)
	return &p, nil
}

func (instance *GameInstance) RemovePlayer(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	delete(instance.Players, username)
	delete(instance.PlayerDataChannels, username)
}

func NewInstance() *GameInstance {
	var gameInstance GameInstance
	gameInstance.Players = make(map[string]*PlayerGameData)
	gameInstance.Rooms = make(map[string]GameRoom)
	gameInstance.InputChannel = make(chan messaging.MessageValue)
	gameInstance.PlayerDataChannels = make(map[string]chan messaging.MessageValue)
	return &gameInstance
}

func (g *GameInstance) GameInstanceRun() {
	for {
		if len(g.Players) != 0 {
			select {
			case val := <-g.InputChannel:
				switch val.GetMessageType() {
				case messaging.MessageOkOrError:
					log.Println("should not be here")
					break
				case messaging.MessageTypeCreateRoom:
					var message *messaging.CommMessageCreateRoom
					message = val.(*messaging.CommMessageCreateRoom)
					if p, ok := g.PlayerDataChannels[message.Player]; ok {
						var okMessage messaging.CommMessageOkOrError
						okMessage.Message = ""
						p <- &okMessage
					}
					break
				case messaging.MessageTypeDeleteRoom:
					var message *messaging.CommMessageDeleteRoom
					message = val.(*messaging.CommMessageDeleteRoom)
					if p, ok := g.PlayerDataChannels[message.Player]; ok {
						var okMessage messaging.CommMessageOkOrError
						okMessage.Message = ""
						p <- &okMessage
					}
					break
				}
				break
			}
		} else {
			time.Sleep(time.Second)
		}

	}
}

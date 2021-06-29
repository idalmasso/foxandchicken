package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameInstance struct {
	Rooms              map[string]*GameRoom
	Players            map[string]struct{}
	PlayersWaiting     map[string]struct{}
	mutex              sync.Mutex
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
	instance.Players[username] = struct{}{}
	instance.PlayersWaiting[username] = struct{}{}
	instance.PlayerDataChannels[username] = make(chan messaging.MessageValue)
	return &p, nil
}

func (instance *GameInstance) RemovePlayer(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	for _, r := range instance.Rooms {
		if _, ok := r.Players[username]; ok {
			r.InputChannel <- &messaging.CommRoomMessageLeftPlayer{Player: username}
		}
	}
	delete(instance.Players, username)
	delete(instance.PlayersWaiting, username)
	delete(instance.PlayerDataChannels, username)
}

func NewInstance() *GameInstance {
	var gameInstance GameInstance
	gameInstance.Players = make(map[string]struct{})
	gameInstance.PlayersWaiting = make(map[string]struct{})
	gameInstance.Rooms = make(map[string]*GameRoom)
	gameInstance.InputChannel = make(chan messaging.MessageValue)
	gameInstance.PlayerDataChannels = make(map[string]chan messaging.MessageValue)
	return &gameInstance
}

func (g *GameInstance) GameInstanceRun() {
	for {
		if len(g.PlayersWaiting) != 0 {
			select {
			case val := <-g.InputChannel:
				switch val.GetMessageType() {
				case messaging.MessageResponse:
					log.Println("should not be here")
				case messaging.MessageTypeCreateRoom:
					var message *messaging.CommMessageCreateRoom
					message = val.(*messaging.CommMessageCreateRoom)
					if p, ok := g.PlayerDataChannels[message.Player]; ok {
						var okMessage messaging.CommMessageResponseCreateRoom
						okMessage.Message = ""

						if _, ok = g.PlayersWaiting[message.Player]; !ok {
							okMessage.Message = "Player already inside a room"
							p <- &okMessage
						} else {
							g.mutex.Lock()
							if _, ok = g.Rooms[message.Name]; !ok {
								room := createRoom(message.Name, g)
								g.Rooms[room.ID] = room
								okMessage.RoomChannel = room.InputChannel
								player := PlayerGameData{Username: message.Player}
								room.Players[message.Player] = &player
								delete(g.PlayersWaiting, message.Player)
								go room.Run()
								p <- &okMessage
							} else {
								okMessage.Message = "Room already exists"
								p <- &okMessage
							}
							g.mutex.Unlock()
						}
					}
				case messaging.RoomMessageTypeJoinPlayer:
					m := val.(*messaging.CommRoomMessageJoinPlayer)
					var r messaging.CommRoomMessageResponse
					if room, ok := g.Rooms[m.Name]; !ok {
						log.Println("Player", m.Player, "try to join not exists room", m.Name)
						r.Message = "room not exists"
						g.mutex.Lock()
						g.PlayerDataChannels[m.Player] <- &r
						g.mutex.Unlock()
					} else {
						log.Println("Player", m.Player, "joined room", room.ID)
						room.mutex.Lock()
						if _, ok := g.Players[m.Player]; ok {
							r.Message = "already exists"
						} else {
							room.broadcastMessage(m)
							p := PlayerGameData{Username: m.Player}
							room.Players[m.Player] = &p
							g.mutex.Lock()
							delete(g.PlayersWaiting, m.Player)
							g.mutex.Unlock()
						}
						room.Instance.PlayerDataChannels[m.Player] <- &r
						room.mutex.Unlock()
					}

				} //end switch messageType
			} //End select
		} else { //if len(g.PlayersWaiting) == 0
			time.Sleep(time.Second)
		}

	}
}

func (g *GameInstance) removeRoom(room string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.Rooms, room)
	var m messaging.CommMessageDeleteRoom
	m.Name = room
	g.broadCastMessageWaitingPlayers(&m)
}

func (g *GameInstance) broadCastMessageWaitingPlayers(message messaging.MessageValue) {
	for p := range g.PlayersWaiting {
		g.PlayerDataChannels[p] <- message
	}
}

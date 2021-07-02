package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameInstance struct {
	Rooms                        map[string]*GameRoom
	Players                      map[string]string //contains the room where he is
	PlayersWaiting               map[string]struct{}
	mutex                        sync.Mutex
	PlayerDataChannels           map[string]chan messaging.InstanceMessageValue
	PlayerDataChannelsBroadcasts map[string]chan messaging.InstanceMessageValue
	InputChannel                 chan messaging.InstanceMessageValue
}

func (instance *GameInstance) AddPlayer(username string) error {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if _, ok := instance.Players[username]; ok {
		return fmt.Errorf("already exists")
	}
	instance.Players[username] = ""
	instance.PlayersWaiting[username] = struct{}{}
	instance.PlayerDataChannels[username] = make(chan messaging.InstanceMessageValue)
	instance.PlayerDataChannelsBroadcasts[username] = make(chan messaging.InstanceMessageValue)
	return nil
}

func (instance *GameInstance) setPlayerWaiting(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	instance.PlayersWaiting[username] = struct{}{}
}
func (instance *GameInstance) RemovePlayer(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if roomName, ok := instance.Players[username]; ok {
		if room, ok := instance.Rooms[roomName]; ok {
			room.RemovePlayer(username)
		}

	}

	delete(instance.Players, username)
	delete(instance.PlayersWaiting, username)
	delete(instance.PlayerDataChannels, username)
	delete(instance.PlayerDataChannelsBroadcasts, username)
}

func NewInstance() *GameInstance {
	var gameInstance GameInstance
	gameInstance.Players = make(map[string]string)
	gameInstance.PlayersWaiting = make(map[string]struct{})
	gameInstance.Rooms = make(map[string]*GameRoom)
	gameInstance.InputChannel = make(chan messaging.InstanceMessageValue)
	gameInstance.PlayerDataChannels = make(map[string]chan messaging.InstanceMessageValue)
	gameInstance.PlayerDataChannelsBroadcasts = make(map[string]chan messaging.InstanceMessageValue)
	return &gameInstance
}

func (g *GameInstance) GameInstanceRun() {
	log.Println("Game instance starting")
	for {
		if len(g.Players) != 0 {
			select {
			case val := <-g.InputChannel:
				switch val.GetMessageType() {
				case messaging.MessageResponse:
					log.Println("should not be here")
				case messaging.MessageTypeCreateRoom:
					message := val.(*messaging.CommMessageCreateRoom)
					if p, ok := g.PlayerDataChannels[message.Player]; ok {
						log.Println("creating room", message.Name)
						var okMessage messaging.CommMessageResponseCreateRoom
						okMessage.Message = ""
						if _, ok = g.PlayersWaiting[message.Player]; !ok {
							okMessage.Message = "Player already inside a room"
							p <- &okMessage
						} else {
							g.mutex.Lock()
							if _, ok = g.Rooms[message.Name]; !ok {
								room := createRoom(message.Name, g)
								g.Rooms[room.Name] = room
								okMessage.RoomChannel = room.RoomInputChannel
								player := PlayerGameData{Username: message.Player}
								room.Players[message.Player] = &player
								room.RoomOutputChannels[message.Player] = make(chan messaging.RoomMessageValue)
								okMessage.RoomResponseChannel = room.RoomOutputChannels[message.Player]
								delete(g.PlayersWaiting, message.Player)
								g.Players[message.Player] = room.Name
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
					message := val.(*messaging.CommRoomMessageJoinPlayer)
					var r messaging.CommMessageResponseJoinRoom
					if room, ok := g.Rooms[message.Name]; !ok {
						log.Println("Player", message.Player, "try to join not existing room", message.Name)
						r.Message = "room not exists"
						g.mutex.Lock()
						g.PlayerDataChannels[message.Player] <- &r
						g.mutex.Unlock()
					} else {
						log.Println("Player", message.Player, "joined room", room.Name)
						room.mutex.Lock()
						if roomExists, ok := g.Players[message.Player]; ok {
							if roomExists == "" {
								room.broadcastMessage(message)
								p := PlayerGameData{Username: message.Player}
								room.Players[message.Player] = &p
								room.RoomOutputChannels[message.Player] = make(chan messaging.RoomMessageValue)
								r.RoomResponseChannel = room.RoomOutputChannels[message.Player]
								r.RoomChannel = room.RoomInputChannel
								g.mutex.Lock()
								delete(g.PlayersWaiting, message.Player)
								g.Players[message.Player] = room.Name
								g.mutex.Unlock()
							} else if roomExists != message.Name {
								r.Message = "already in room " + roomExists
							}
						}
						g.PlayerDataChannels[message.Player] <- &r
						room.mutex.Unlock()
					}
				case messaging.RoomMessageTypeLeftPlayer:
					message := val.(*messaging.CommRoomMessageLeftPlayer)
					log.Println("Player", message.Player, "left room", g.Players[message.Player])

					if message.Player == "" {
						log.Println("empty player message")

					} else {
						response := messaging.CommMessageResponse{Message: ""}
						if _, ok := g.Rooms[g.Players[message.Player]]; ok {
							log.Println("ok sending message to player back")
							g.setPlayerWaiting(message.Player)
							g.Rooms[g.Players[message.Player]].RemovePlayer(message.Player)
							g.PlayerDataChannels[message.Player] <- &response
						} else {
							response.Message = "not found room"
						}
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
	g.broadCastMessageAllPlayers(&m)
}

func (g *GameInstance) broadCastMessageWaitingPlayers(message messaging.InstanceMessageValue) {
	for p := range g.PlayersWaiting {
		g.PlayerDataChannelsBroadcasts[p] <- message
	}
}

func (g *GameInstance) broadCastMessageAllPlayers(message messaging.InstanceMessageValue) {
	for p := range g.Players {
		g.PlayerDataChannelsBroadcasts[p] <- message
	}
}

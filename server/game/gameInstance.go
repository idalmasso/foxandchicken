package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

//GameInstance struct containing the data for an instance and its rooms
type GameInstance struct {
	Rooms                        map[string]*GameRoom
	Players                      map[string]string //contains the room where he is
	PlayersWaiting               map[string]struct{}
	mutex                        sync.Mutex
	PlayerDataChannels           map[string]chan messaging.InstanceMessageValue
	PlayerDataChannelsBroadcasts map[string]chan messaging.InstanceMessageValue
	InputChannel                 chan messaging.InstanceMessageValue
}

//AddPlayer adds a player to the instance
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

//setPlayerWaiting set a player as waiting (not in a room)
func (instance *GameInstance) setPlayerWaiting(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	instance.PlayersWaiting[username] = struct{}{}
}

//RemovePlayer removes a player from the instance
func (instance *GameInstance) RemovePlayer(username string) {

	if roomName, ok := instance.Players[username]; ok {
		if room, ok := instance.Rooms[roomName]; ok {
			instance.Players[username] = ""
			room.RemovePlayer(username)
		}

	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	delete(instance.Players, username)
	delete(instance.PlayersWaiting, username)
	close(instance.PlayerDataChannels[username])
	delete(instance.PlayerDataChannels, username)
	close(instance.PlayerDataChannelsBroadcasts[username])
	delete(instance.PlayerDataChannelsBroadcasts, username)
}

//NewInstance return a GameInstance
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

//GameInstanceRun is the main instance (creates and remove rooms and other)
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
							if room, err := g.tryCreateRoom(message.Name); err != nil {
								okMessage.Message = err.Error()
							} else {
								g.addPlayerToRoom(room, message.Player)
								okMessage.RoomResponseChannel = room.RoomOutputChannels[message.Player]
								okMessage.RoomChannel = room.RoomInputChannel
								go room.Run()
							}
							p <- &okMessage
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
						if roomOfPlayer, ok := g.Players[message.Player]; ok {
							if roomOfPlayer == "" {
								room.broadcastMessage(message)
								g.addPlayerToRoom(room, message.Player)
								r.RoomResponseChannel = room.RoomOutputChannels[message.Player]
								r.RoomChannel = room.RoomInputChannel
							} else if roomOfPlayer != message.Name {
								r.Message = "already in room " + roomOfPlayer
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
func (g *GameInstance) tryCreateRoom(room string) (*GameRoom, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	log.Println("tryCreateRoom - Start")
	if _, ok := g.Rooms[room]; !ok {
		room := createRoom(room, g)
		g.Rooms[room.Name] = room
		return room, nil
	} else {
		return nil, fmt.Errorf("Room already exists")

	}
}
func (g *GameInstance) addPlayerToRoom(room *GameRoom, player string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	room.AddPlayer(player)
	delete(g.PlayersWaiting, player)
	g.Players[player] = room.Name
}

//removeRoom removes a room from an instance
func (g *GameInstance) removeRoom(room string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.Rooms, room)
	var m messaging.CommMessageDeleteRoom
	m.Name = room
	g.broadCastMessageAllPlayers(&m)
}

//broadCastMessageWaitingPlayers broadcast a message to players not in a room
func (g *GameInstance) broadCastMessageWaitingPlayers(message messaging.InstanceMessageValue) {
	for p := range g.PlayersWaiting {
		g.PlayerDataChannelsBroadcasts[p] <- message
	}
}

//broadCastMessageAllPlayers broadcast a message to all players
func (g *GameInstance) broadCastMessageAllPlayers(message messaging.InstanceMessageValue) {
	for p := range g.Players {
		g.PlayerDataChannelsBroadcasts[p] <- message
	}
}

func (g *GameInstance) GetRooms() []GameRoomNumPlayer {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	rooms := make([]GameRoomNumPlayer, len(g.Rooms))
	counter := 0
	for name, room := range g.Rooms {
		gameRoom := GameRoomNumPlayer{Name: name, Players: len(room.Players)}
		rooms[counter] = gameRoom
		counter++
	}

	return rooms
}

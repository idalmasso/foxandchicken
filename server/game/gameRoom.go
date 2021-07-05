package game

import (
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

//GameRoom struct containing the game room data
type GameRoom struct {
	Name               string `json:"name"`
	Players            map[string]*PlayerGameData
	sizeX, sizeY       float32
	status             int
	Instance           *GameInstance
	mutex              sync.Mutex
	RoomInputChannel   chan messaging.RoomMessageValue
	RoomOutputChannels map[string]chan messaging.RoomMessageValue
	timestamp          int64
}

//createRoom creates the actual room in a gameinstance
func createRoom(name string, instance *GameInstance) *GameRoom {
	log.Println("Creating room", name)
	g := GameRoom{Instance: instance}
	g.status = 0
	g.Players = make(map[string]*PlayerGameData)
	g.RoomInputChannel = make(chan messaging.RoomMessageValue)
	g.Name = name
	g.sizeX, g.sizeY = 100, 100
	g.RoomOutputChannels = make(map[string]chan messaging.RoomMessageValue, 0)
	return &g
}

//Run is the GameRoom main call
func (g *GameRoom) Run() {
	g.mutex.Lock()
	g.timestamp = time.Now().UnixNano()
	g.mutex.Unlock()
	for {
		if len(g.Players) == 0 {
			log.Println("Room", g.Name, "empty, removing")
			g.Instance.removeRoom(g.Name)
			return
		}
		select {
		case val := <-g.RoomInputChannel:
			if val != nil {
				switch val.GetMessageType() {

				case messaging.RoomMessageTypeMovePlayer:
					m := val.(*messaging.CommRoomMessageMovePlayer)
					//todo: Remove this line
					log.Println("Player", m.Player, "move in room", g.Name)
					g.movePlayer(m)

				}
			}
		default:
			g.gameCycle()
		}

	}
}

//broadcastMessage send a message to all players in room
func (g *GameRoom) broadcastMessage(message messaging.RoomMessageValue) {
	//log.Printf("room %s broadcast %T", g.Name, message.GetMessageType())

	for p := range g.Players {
		//log.Println("---Send message to", p)
		g.RoomOutputChannels[p] <- message
	}
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) gameCycle() {
	newTimestamp := time.Now().UnixNano()
	//timeDelta := newTimestamp - g.timestamp
	g.mutex.Lock()
	defer func() {
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 50)
	}()
	message := make(messaging.CommRoomMessagePlayersMovement, len(g.Players))
	i := 0

	for username, p := range g.Players {
		//CALL MOVEPLAYER FOR ALL PLAYERS HERE
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.Position
		m.Rotation = p.Rotation
		m.Velocity = p.Velocity
		m.Player = username
		message[i] = m
		i++
	}
	g.broadcastMessage(&message)
	g.timestamp = newTimestamp

}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) movePlayer(m *messaging.CommRoomMessageMovePlayer) {
	g.Players[m.Player].Position = m.Position
	g.Players[m.Player].Rotation = m.Rotation
	g.Players[m.Player].Velocity = m.Velocity
}

//RemovePlayer removes a player from the room
func (g *GameRoom) RemovePlayer(username string) {
	log.Println("Removing player ", username)
	g.mutex.Lock()
	delete(g.Players, username)
	delete(g.RoomOutputChannels, username)
	g.mutex.Unlock()
	g.broadcastMessage(&messaging.CommRoomMessageLeftPlayer{Player: username})
}

//AddPlayer add a player in the room
func (g *GameRoom) AddPlayer(username string) {
	player := PlayerGameData{Username: username}
	g.Players[username] = &player
	g.RoomOutputChannels[username] = make(chan messaging.RoomMessageValue)
}

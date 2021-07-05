package game

import (
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameRoom struct {
	Name               string `json:"name"`
	Players            map[string]*PlayerGameData
	sizeX, sizeY       float32
	status             int
	Instance           *GameInstance
	mutex              sync.Mutex
	RoomInputChannel   chan messaging.RoomMessageValue
	RoomOutputChannels map[string]chan messaging.RoomMessageValue
}

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

func (g *GameRoom) Run() {
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

func (g *GameRoom) broadcastMessage(message messaging.RoomMessageValue) {
	//log.Printf("room %s broadcast %T", g.Name, message.GetMessageType())
	g.Instance.mutex.Lock()
	defer g.Instance.mutex.Unlock()
	for p := range g.Players {
		//log.Println("---Send message to", p)
		g.RoomOutputChannels[p] <- message
	}
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) gameCycle() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	message := make(messaging.CommRoomMessagePlayersMovement, len(g.Players))
	i := 0

	for username, p := range g.Players {
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.Position
		m.Rotation = p.Rotation
		m.Velocity = p.Velocity
		m.Player = username
		message[i] = m
		i++
	}
	g.broadcastMessage(&message)
	time.Sleep(time.Millisecond * 20)
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) movePlayer(m *messaging.CommRoomMessageMovePlayer) {
	g.Players[m.Player].Position = m.Position
	g.Players[m.Player].Rotation = m.Rotation
	g.Players[m.Player].Velocity = m.Velocity
}

func (g *GameRoom) RemovePlayer(username string) {
	log.Println("Removing player ", username)
	g.mutex.Lock()
	delete(g.Players, username)
	delete(g.RoomOutputChannels, username)
	g.mutex.Unlock()
	//g.broadcastMessage(m)
}

func (g *GameRoom) AddPlayer(username string) {
	player := PlayerGameData{Username: username}
	g.Players[username] = &player
	g.RoomOutputChannels[username] = make(chan messaging.RoomMessageValue)
}

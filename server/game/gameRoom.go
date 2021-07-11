package game

import (
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

//GameRoomNumPlayer is the struct that returns the data about number rooms and players in that
type GameRoomNumPlayer struct{
	Name string `json:"name"`
	Players int `json:"players"`
}
//GameRoom struct containing the game room data
type GameRoom struct {
	Name               string `json:"name"`
	Players            map[string]*PlayerGameData
	sizeX, sizeY       float64
	status             int
	Instance           *GameInstance
	mutex              sync.Mutex
	RoomInputChannel   chan messaging.RoomMessageValue
	RoomOutputChannels map[string]chan messaging.RoomMessageValue
	MaxAcceleration    float64
	MaxVelocity        float64
	Drag               float64
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
	g.MaxAcceleration = 1
	g.MaxVelocity = 2
	g.Drag = 0.9
	g.RoomOutputChannels = make(map[string]chan messaging.RoomMessageValue )
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
					g.playerInput(m)

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
	timeDelta := time.Duration(newTimestamp - g.timestamp)
	g.mutex.Lock()
	defer func() {
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 50)
	}()
	message := make(messaging.CommRoomMessagePlayersMovement, len(g.Players))
	i := 0

	for username, p := range g.Players {
		//CALL MOVEPLAYER FOR ALL PLAYERS HERE
		p.mutex.Lock()
		g.movePlayer(p, timeDelta)
		p.timestamp = newTimestamp
		p.mutex.Unlock()
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.Position
		m.Rotation = p.Rotation
		m.Velocity = p.Velocity
		m.Player = username
		m.Timestamp = newTimestamp
		message[i] = m
		i++
	}
	g.broadcastMessage(&message)
	g.timestamp = newTimestamp

}

func (g *GameRoom) movePlayer(p *PlayerGameData, deltaT time.Duration) {
	ts := deltaT.Seconds()
	p.Position = common.VectorSum(p.Position, p.Velocity.ScalarProduct(ts))
	p.Position = p.Position.ClampVector(0, g.sizeX, 0, g.sizeY)
	if p.Position.X < 0 {
		p.Position.X = 0
	}
	if p.Position.X > g.sizeX {
		p.Position.X = g.sizeX
	}
	if p.Acceleration.X == 0 && p.Acceleration.Y == 0 {
		magnitude := p.Velocity.SqrtMagnitude()
		if magnitude < 0.01 {
			p.Velocity.X = 0
			p.Velocity.Y = 0
			return
		}
		p.Velocity = common.VectorSum(p.Velocity, p.Velocity.ScalarProduct(-g.Drag*ts))
	} else {
		p.Velocity = common.VectorSum(p.Velocity, p.Acceleration.ScalarProduct(ts))
		magnitude := p.Velocity.SqrtMagnitude()

		if magnitude > float64(g.MaxVelocity) {
			p.Velocity = p.Velocity.ScalarProduct(g.MaxVelocity / magnitude)
		}
	}

}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) playerInput(m *messaging.CommRoomMessageMovePlayer) {
	g.Players[m.Player].mutex.Lock()
	defer g.Players[m.Player].mutex.Unlock()
	newTimestamp := time.Now().UnixNano()
	if m.Timestamp > newTimestamp || m.Timestamp == 0 {
		m.Timestamp = newTimestamp
	}

	magnitude := m.Acceleration.SqrtMagnitude()

	if magnitude > float64(g.MaxAcceleration) {
		m.Acceleration = m.Acceleration.ScalarProduct(g.MaxAcceleration / magnitude)
	}
	g.Players[m.Player].Rotation = m.Rotation
	g.Players[m.Player].Acceleration = m.Acceleration
	g.movePlayer(g.Players[m.Player], time.Duration(newTimestamp-m.Timestamp))
	g.timestamp = newTimestamp

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
	player := PlayerGameData{Username: username, timestamp: time.Now().UnixNano()}
	g.Players[username] = &player
	g.RoomOutputChannels[username] = make(chan messaging.RoomMessageValue)
}

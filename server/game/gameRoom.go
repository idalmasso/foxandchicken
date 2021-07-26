package game

import (
	"log"
	"sync"
	"time"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

//GameRoomNumPlayer is the struct that returns the data about number rooms and players in that
type GameRoomNumPlayer struct {
	Name    string `json:"name"`
	Players int    `json:"players"`
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
	g.sizeX, g.sizeY = 20, 20
	g.MaxAcceleration = 8
	g.MaxVelocity = 2
	g.Drag = 0.95
	g.RoomOutputChannels = make(map[string]chan messaging.RoomMessageValue)
	return &g
}

//Run is the GameRoom main call
func (g *GameRoom) Run() {
	log.Println("Run - Lock")
	g.mutex.Lock()
	g.timestamp = time.Now().UnixNano()
	log.Println("Run - UnLock")
	g.mutex.Unlock()
	for {
		if len(g.Players) == 0 {
			log.Println("Room", g.Name, "empty, removing")
			g.Instance.removeRoom(g.Name)
			return
		}
		select {
		case val := <-g.RoomInputChannel:
			log.Println("Read room input channel")
			if val != nil {
				switch val.GetMessageType() {
				case messaging.RoomMessageTypeMovePlayer:
					m := val.(*messaging.CommRoomMessageMovePlayer)
					g.playerInput(m)

				}
			} else {
				log.Println("Got a null room message")
			}
		default:
			log.Println("Game cycle")
			g.gameCycle()
			log.Println("End Game cycle")
		}

	}
}

//broadcastMessage send a message to all players in room
func (g *GameRoom) broadcastMessage(message messaging.RoomMessageValue) {
	log.Printf("room %s broadcast %T", g.Name, message.GetMessageType())

	for p := range g.Players {
		//log.Println("---Send message to", p)
		if g.Instance.Players[p] == g.Name {
			g.RoomOutputChannels[p] <- message
		}
	}
	log.Printf("room %s broadcasted %T", g.Name, message.GetMessageType())
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) gameCycle() {
	newTimestamp := time.Now().UnixNano()
	log.Println("gameCycle - Lock")
	g.mutex.Lock()
	defer func() {
		log.Println("gameCycle - Unlock")
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 50)
	}()
	log.Println("gameCycle - BEFORE MAKE")
	message := make(messaging.CommRoomMessagePlayersMovement, len(g.Players))
	i := 0
	deltaT := time.Duration(newTimestamp - g.timestamp).Seconds()
	//Will be range of gameobjects
	for username, p := range g.Players {
		log.Println("gameCycle - BEFORE LOCK")
		p.mutex.Lock()
		log.Println("gameCycle - AFTER Lock")
		p.gameObject.Update(deltaT)
		log.Println("gameCycle - AFTER Update")
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.gameObject.Position

		p.mutex.Unlock()
		log.Println("gameCycle - AFTER uNLock")
		m.Player = username
		m.Timestamp = newTimestamp
		message[i] = m
		i++
	}
	log.Println("Game cycle broadcasting move")
	log.Println("gameCycle - before broadcastMessage")
	g.broadcastMessage(&message)
	log.Println("gameCycle - after broadcastMessage")
	g.timestamp = newTimestamp

}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) playerInput(m *messaging.CommRoomMessageMovePlayer) {
	magnitude := m.Acceleration.SqrtMagnitude()
	if magnitude != 0 {
		m.Acceleration = m.Acceleration.ScalarProduct(g.MaxAcceleration / magnitude)
	}
	g.Players[m.Player].SetInput(m.Acceleration)
}

//RemovePlayer removes a player from the room
func (g *GameRoom) RemovePlayer(username string) {
	log.Println("Room removing player ", username)
	log.Println("RemovePlayer - Lock")
	g.mutex.Lock()
	defer func() {
		log.Println("RemovePlayer - Unlock")
		g.mutex.Unlock()
	}()

	delete(g.Players, username)
	delete(g.RoomOutputChannels, username)

	g.broadcastMessage(&messaging.CommRoomMessageLeftPlayer{Player: username})
}

//AddPlayer add a player in the room
func (g *GameRoom) AddPlayer(username string) {
	g.Players[username] = NewPlayer(username, 1, g)
	g.RoomOutputChannels[username] = make(chan messaging.RoomMessageValue)

}

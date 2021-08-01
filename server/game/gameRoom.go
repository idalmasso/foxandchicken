package game

import (
	"sync"
	"time"

	"github.com/golang/glog"
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
	roomStopChannel    chan bool
	cellsToGameObjectmap  map[int]map[*GameObject]struct{}
	
}

//createRoom creates the actual room in a gameinstance
func createRoom(name string, instance *GameInstance) *GameRoom {
	if glog.V(2) {
		glog.Infoln("creating room", name)
	}
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
	g.cellsToGameObjectmap = make(map[int]map[*GameObject]struct{})
	for x:=0; x<int(g.sizeX*g.sizeY); x++{
		g.cellsToGameObjectmap[x]= make(map[*GameObject]struct{})
	}
	g.roomStopChannel = make(chan bool)
	return &g
}

//Run is the GameRoom main call
func (g *GameRoom) Run() {
	if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.Run - Lock")
	}
	g.mutex.Lock()
	g.timestamp = time.Now().UnixNano()
	if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.Run - UnLock")
	}
	g.mutex.Unlock()
	go g.gameCycle()
	for {
		select {
		case <-g.roomStopChannel:
			if glog.V(3) {
				glog.Infoln("DEBUG - GameRoom.Run - Room", g.Name, "Game run stop")
			}
			return
		case val := <-g.RoomInputChannel:
			if glog.V(3) {
				glog.Infoln("DEBUG - GameRoom.Run - Read room input channel")
			}
			if val != nil {
				switch val.GetMessageType() {
				case messaging.RoomMessageTypeMovePlayer:
					m := val.(*messaging.CommRoomMessageMovePlayer)
					g.playerInput(m)
				}
				
			} else {
				if glog.V(1) {
					glog.Warningln("GameRoom.Run - Got a null room message")
				}
			}

		}

	}
}

//broadcastMessage send a message to all players in room
func (g *GameRoom) broadcastMessage(message messaging.RoomMessageValue) {
	if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.broadcastMessage - room", g.Name, "broadcast", message.GetMessageType())
	}

	for p := range g.Players {
		//log.Println("---Send message to", p)
		if g.Instance.Players[p] == g.Name {
			g.RoomOutputChannels[p] <- message
		}
	}
	if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.broadcastMessage - room", g.Name, "broadcasted", message.GetMessageType())
	}
}

//gameCycle is the output cycle, and also the recalculate. Probably should do in 2 different goroutines, to be more clean (1 for the "game", one for the "output")
func (g *GameRoom) gameCycle() {
	for {
		select {
		case <-g.roomStopChannel:
			if glog.V(2) {
				glog.Infoln("GameRoom.gameCycle - Room", g.Name, "Game cycle stop")
			}
			return
		default:
			g.updateAndSendData()
		}
	}
}
func (g *GameRoom) updateAndSendData() {

	g.mutex.Lock()
	defer func() {
		g.mutex.Unlock()
		time.Sleep(time.Millisecond * 50)
	}()
	if len(g.Players) <= 0 {
		return
	}
	newTimestamp := time.Now().UnixNano()
	message := make(messaging.CommRoomMessagePlayersMovement, len(g.Players))
	i := 0
	deltaT := time.Duration(newTimestamp - g.timestamp).Seconds()
	//Will be range of gameobjects
	for username, p := range g.Players {
		p.mutex.Lock()
		p.gameObject.Update(deltaT)
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.gameObject.Position

		p.mutex.Unlock()
		m.Player = username
		m.Timestamp = newTimestamp
		message[i] = m
		i++
	}
	g.broadcastMessage(&message)
	g.timestamp = newTimestamp
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) playerInput(m *messaging.CommRoomMessageMovePlayer) {
	magnitude := m.Acceleration.SqrtMagnitude()
	if magnitude != 0 {
		m.Acceleration = m.Acceleration.ScalarProduct(g.MaxAcceleration / magnitude)
	}
	g.Players[m.Player].SetInput(m.Acceleration, m.ActionPressed)
}

//RemovePlayer removes a player from the room
func (g *GameRoom) RemovePlayer(username string) {
	if glog.V(2) {
		glog.Infoln("GameRoom.RemovePlayer - removing player ", username)
	}
	if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.RemovePlayer - Lock")
	}
	g.mutex.Lock()
	defer func() {
		if glog.V(3) {
		glog.Infoln("DEBUG - GameRoom.RemovePlayer - UnLock")
	}
		g.mutex.Unlock()
	}()

	delete(g.Players, username)
	delete(g.RoomOutputChannels, username)

	g.broadcastMessage(&messaging.CommRoomMessageLeftPlayer{Player: username})
	if len(g.Players) <= 0 {
		g.Instance.removeRoom(g.Name)
		g.roomStopChannel <- true //One for the input goroutine
		g.roomStopChannel <- true //One for the output goroutine
		close(g.roomStopChannel)
	}
}

//AddPlayer add a player in the room
func (g *GameRoom) AddPlayer(username string) {
	if len(g.Players)%2==0{
		g.Players[username] = NewPlayer(username, CharacterTypeFox, g)
	} else {
		g.Players[username] = NewPlayer(username, CharacterTypeChicken, g)
	}
	g.RoomOutputChannels[username] = make(chan messaging.RoomMessageValue)
	
	g.addGameObject(g.Players[username].gameObject)
}

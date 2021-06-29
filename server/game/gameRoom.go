package game

import (
	"log"
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameRoom struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Players      map[string]*PlayerGameData
	sizeX, sizeY float32
	status       int
	Instance     *GameInstance
	mutex        sync.Mutex
	InputChannel chan messaging.RoomMessageValue
}

func createRoom(name string, instance *GameInstance) *GameRoom {
	log.Println("Creating room", name)
	g := GameRoom{Instance: instance}
	g.status = 0
	g.Players = make(map[string]*PlayerGameData)
	g.InputChannel = make(chan messaging.RoomMessageValue)
	g.Name = name
	g.ID = name
	g.sizeX, g.sizeY = 100, 100
	return &g
}

func (g *GameRoom) Run() {
	for {
		if len(g.Players) == 0 {
			log.Println("Room", g.ID, "empty, removing")
			g.Instance.removeRoom(g.ID)
			return
		}
		select {
		case val := <-g.InputChannel:
			switch val.GetMessageType() {

			case messaging.RoomMessageTypeLeftPlayer:
				m := val.(*messaging.CommRoomMessageLeftPlayer)
				log.Println("Player", m.Player, "left room", g.ID)
				g.mutex.Lock()
				var r messaging.CommRoomMessageResponse
				delete(g.Players, m.Player)
				g.Instance.PlayerDataChannels[m.Player] <- &r
				g.mutex.Unlock()
				g.broadcastMessage(m)
			case messaging.RoomMessageTypeMovePlayer:
				m := val.(*messaging.CommRoomMessageMovePlayer)
				//todo: Remove this line
				log.Println("Player", m.Player, "move in room", g.ID)
				g.movePlayer(m)
			default:
				g.gameCycle()
			}
		}

	}
}

func (g *GameRoom) broadcastMessage(message messaging.MessageValue) {
	log.Println("room", g.ID, "broadcast", message.GetMessageType())
	g.Instance.mutex.Lock()
	defer g.Instance.mutex.Unlock()
	for p := range g.Players {
		g.Instance.PlayerDataChannels[p] <- message
	}
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) gameCycle() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, p := range g.Players {
		var m messaging.CommRoomMessageMovePlayer
		m.Position = p.Position
		m.Rotation = p.Rotation
		m.Velocity = p.Velocity
		g.broadcastMessage(&m)
	}
}

//Right by now it will be ALL on frontend... Next->checks
func (g *GameRoom) movePlayer(m *messaging.CommRoomMessageMovePlayer) {
	g.Players[m.Player].Position = m.Position
	g.Players[m.Player].Rotation = m.Rotation
	g.Players[m.Player].Velocity = m.Velocity
}

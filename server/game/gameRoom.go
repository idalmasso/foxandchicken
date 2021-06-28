package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/messaging"
)

type GameRoom struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Players      map[string]struct{}
	status       int
	Instance     *GameInstance
	mutex        sync.Mutex
	InputChannel chan messaging.RoomMessageValue
}

func createRoom(name string, instance *GameInstance) *GameRoom {
	g := GameRoom{Instance: instance}
	g.status = 0
	g.Players = make(map[string]struct{})
	g.InputChannel = make(chan messaging.RoomMessageValue)
	g.Name = name
	g.ID = name
	return &g
}

func (g *GameRoom) Run() {
	for {
		if len(g.Players) == 0 {
			g.Instance.removeRoom(g.ID)
			return
		}
		select {
		case val := <-g.InputChannel:
			switch val.GetMessageType() {
			case messaging.RoomMessageTypeJoinPlayer:
				m := val.(*messaging.CommRoomMessageJoinPlayer)
				var r messaging.CommRoomMessageResponse
				g.mutex.Lock()
				if _, ok := g.Players[m.Player]; ok {
					r.Message = "already exists"
				} else {
					g.broadcastMessage(m)
					g.Players[m.Player] = struct{}{}
				}
				g.Instance.PlayerDataChannels[m.Player] <- &r
				g.mutex.Unlock()
				break
			case messaging.RoomMessageTypeLeftPlayer:
				m := val.(*messaging.CommRoomMessageLeftPlayer)
				g.mutex.Lock()
				var r messaging.CommRoomMessageResponse
				delete(g.Players, m.Player)
				g.Instance.PlayerDataChannels[m.Player] <- &r
				g.mutex.Unlock()
				g.broadcastMessage(m)
				break
			}
		}
	}
}

func (g *GameRoom) broadcastMessage(message messaging.MessageValue) {
	g.Instance.mutex.Lock()
	defer g.Instance.mutex.Unlock()
	for p, _ := range g.Players {
		g.Instance.PlayerDataChannels[p] <- message
	}
}

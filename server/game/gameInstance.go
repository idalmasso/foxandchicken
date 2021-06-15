package game

import (
	"fmt"
	"sync"
)

type GameInstance struct {
	Rooms   map[string]GameRoom
	Players map[string]*PlayerGameData
	mutex   sync.Mutex
}

func (instance *GameInstance) AddPlayer(username string) (*PlayerGameData, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if _, ok := instance.Players[username]; ok {
		return nil, fmt.Errorf("already exists")
	}
	var p PlayerGameData
	p.Username = username
	instance.Players[username] = &p
	return &p, nil
}

func (instance *GameInstance) RemovePlayer(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	delete(instance.Players, username)
}

func NewInstance() *GameInstance {
	var gameInstance GameInstance
	gameInstance.Players = make(map[string]*PlayerGameData)
	gameInstance.Rooms = make(map[string]GameRoom)
	return &gameInstance
}

package game

import (
	"fmt"
	"sync"
)

type GameInstance struct {
	Rooms   map[string]GameRoom
	Players map[string]bool
	mutex   sync.Mutex
}

func (instance *GameInstance) AddPlayer(username string) (string, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if _, ok := instance.Players[username]; ok {
		return "", fmt.Errorf("already exists")
	}
	instance.Players[username] = true
	return username, nil
}

func (instance *GameInstance) RemovePlayer(username string) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	delete(instance.Players, username)
}

func NewInstance() *GameInstance {
	var gameInstance GameInstance
	gameInstance.Players = make(map[string]bool)
	gameInstance.Rooms = make(map[string]GameRoom)
	return &gameInstance
}

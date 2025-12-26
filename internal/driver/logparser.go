package driver

import (
	"maps"
	"slices"
	"strings"
	"sync"
)

type OnlinePlayers struct {
	mutex   sync.RWMutex
	players map[string]struct{}
	drv     *Driver
}

func NewOnlinePlayers(driver *Driver) *OnlinePlayers {
	return &OnlinePlayers{
		mutex:   sync.RWMutex{},
		players: make(map[string]struct{}),
		drv:     driver,
	}
}

func (op *OnlinePlayers) Get() []string {
	op.mutex.RLock()
	defer op.mutex.RUnlock()

	result := slices.Collect(maps.Keys(op.players))
	if result == nil {
		return []string{}
	}
	return result
}

func (op *OnlinePlayers) Count() int {
	op.mutex.RLock()
	defer op.mutex.RUnlock()
	return len(op.players)
}

func (op *OnlinePlayers) parseLine(line string) {
	words := strings.Split(line[:len(line)-1], " ")
	if len(words) != 7 {
		return
	}
	playerName := words[3]
	didChange := false

	if words[5] == "the" && words[6] == "game" && playerName[0] != '<' {
		switch words[4] {
		case "joined":
			op.mutex.Lock()
			op.players[playerName] = struct{}{}
			didChange = true
			op.mutex.Unlock()
		case "left":
			op.mutex.Lock()
			delete(op.players, playerName)
			didChange = true
			op.mutex.Unlock()
		}
	}

	if didChange {
		op.drv.stateBroadcaster.sendUpdate(op.drv.GetState())
		if op.Count() == 0 {
			op.drv.ScheduleStop()
		}
	}
}

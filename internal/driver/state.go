package driver

import (
	"sync"
)

type ServerState struct {
	Version       string   `json:"version"`
	ConnectUrl    string   `json:"connect_url"`
	ServerName    string   `json:"server_name"`
	IsRunning     bool     `json:"is_running"`
	OnlinePlayers []string `json:"online_players"` // TODO
	AutoStopTime  int64    `json:"auto_stop_time"` // TODO
	BotTag        string   `json:"bot_tag"`        // TODO
}

type StateBroadcaster struct {
	updatesSubscribers      map[uint64]func(state *ServerState)
	updatesSubscribersMutex sync.RWMutex
	driver                  *Driver
}

func NewStateBroadcaster(drv *Driver) *StateBroadcaster {
	return &StateBroadcaster{
		updatesSubscribers:      make(map[uint64]func(state *ServerState)),
		updatesSubscribersMutex: sync.RWMutex{},
		driver:                  drv,
	}
}

func (sb *StateBroadcaster) Subscribe(id uint64, updater func(*ServerState)) {
	sb.updatesSubscribersMutex.Lock()
	sb.updatesSubscribers[id] = updater
	sb.updatesSubscribersMutex.Unlock()
}

func (sb *StateBroadcaster) Unsubscribe(id uint64) {
	sb.updatesSubscribersMutex.Lock()
	delete(sb.updatesSubscribers, id)
	sb.updatesSubscribersMutex.Unlock()
}

func (sb *StateBroadcaster) sendUpdate() {
	state := sb.driver.GetState()
	sb.updatesSubscribersMutex.RLock()
	for _, updater := range sb.updatesSubscribers {
		go func(updater func(*ServerState)) {
			updater(state)
		}(updater)
	}
	sb.updatesSubscribersMutex.RUnlock()
}

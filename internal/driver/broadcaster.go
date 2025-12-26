package driver

import (
	"sync"
)

type Broadcaster[UpdateT any] struct {
	subscribers      map[uint64]func(update UpdateT)
	subscriversMutex sync.RWMutex
}

func NewBroadcaster[UpdateT any]() *Broadcaster[UpdateT] {
	return &Broadcaster[UpdateT]{
		subscribers:      make(map[uint64]func(update UpdateT)),
		subscriversMutex: sync.RWMutex{},
	}
}

func (b *Broadcaster[UpdateT]) Subscribe(id uint64, updater func(UpdateT)) {
	b.subscriversMutex.Lock()
	b.subscribers[id] = updater
	b.subscriversMutex.Unlock()
}

func (b *Broadcaster[UpdateT]) Unsubscribe(id uint64) {
	b.subscriversMutex.Lock()
	delete(b.subscribers, id)
	b.subscriversMutex.Unlock()
}

func (b *Broadcaster[UpdateT]) sendUpdate(update UpdateT) {
	b.subscriversMutex.RLock()
	for _, updater := range b.subscribers {
		go func(updater func(UpdateT)) {
			updater(update)
		}(updater)
	}
	b.subscriversMutex.RUnlock()
}

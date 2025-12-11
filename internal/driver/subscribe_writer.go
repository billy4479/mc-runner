package driver

import (
	"io"
	"sync"

	"go.uber.org/multierr"
)

type SubscribeWriter struct {
	subscribers map[uint64]io.Writer
	wg          sync.WaitGroup
	mutex       sync.RWMutex // mutex for the map, not the writers
}

func NewSubscribeWriter() *SubscribeWriter {
	return &SubscribeWriter{
		subscribers: make(map[uint64]io.Writer),
		wg:          sync.WaitGroup{},
		mutex:       sync.RWMutex{},
	}
}

func (sr *SubscribeWriter) Subscribe(w io.Writer, id uint64) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	sr.subscribers[id] = w
}

func (sr *SubscribeWriter) Unsubscribe(id uint64) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	delete(sr.subscribers, id)
}

func (sr *SubscribeWriter) Write(b []byte) (int, error) {
	sr.mutex.RLock()

	nMin := len(b)
	var errs error
	m := sync.Mutex{} // for errs

	for _, dst := range sr.subscribers {
		sr.wg.Add(1)
		go func(dst io.Writer) {
			n, err := dst.Write(b)

			nMin = min(nMin, n)
			if err != nil {
				m.Lock()
				errs = multierr.Append(errs, err)
				m.Unlock()
			}

			sr.wg.Done()
		}(dst)
	}

	sr.mutex.RUnlock()

	sr.wg.Wait()

	return nMin, errs
}

package internal

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type State struct {
	ws      map[uint64]*websocket.Conn
	wsMutex sync.Mutex

	serverLog       string
	isServerRunning bool
	runningMutex    sync.Mutex
}

func NewState() *State {
	return &State{
		ws:              make(map[uint64]*websocket.Conn),
		wsMutex:         sync.Mutex{},
		serverLog:       "",
		isServerRunning: false,
		runningMutex:    sync.Mutex{}}
}

func (wss *State) AddConnection(ws *websocket.Conn, id uint64) {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	wss.ws[id] = ws
}

func (wss *State) GetConnection(id uint64) *websocket.Conn {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	return wss.ws[id]
}

func (wss *State) CloseAndRemoveConnection(id uint64) error {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	err := wss.ws[id].WriteMessage(websocket.CloseNormalClosure, nil)
	delete(wss.ws, id)

	log.Debug().Uint64("request_id", id).Msg("ws closed")

	return err
}

func addWebsocket(g *echo.Group, config *Config) {
	wss := NewState()

	idMutex := sync.Mutex{}
	reqId := uint64(0)

	g.GET("/ws", func(c echo.Context) error {
		idMutex.Lock()
		connId := reqId
		reqId++
		idMutex.Unlock()

		logIfErr := func(err error) bool {
			if err == nil {
				return false
			}
			log.Error().Err(fmt.Errorf("websocket %d: %w", connId, err))
			return true
		}
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		wss.AddConnection(ws, connId)
		defer logIfErr(wss.CloseAndRemoveConnection(connId))

		log.Info().Uint64("request_id", connId).Msg("connection upgraded")

		for {
			msgType, msg, err := ws.ReadMessage()
			if logIfErr(err) {
				break
			}

			if msgType != websocket.TextMessage {
				log.Warn().Uint64("request_id", connId).Int("msg_type", msgType).Msg("got non-text message")
				continue
			}

			switch string(msg) {
			case "ping":
				err := ws.WriteJSON(echo.Map{"type": "version", "data": "mc-runner@" + Version})
				logIfErr(err)

			case "start":
				log.Info().Msg("TODO: starting server")
			}
		}
		return nil
	}, authMiddleware)
}

package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/billy4479/mc-runner/internal/driver"
)

type WsState struct {
	ws      map[uint64]*websocket.Conn
	wsMutex sync.Mutex

	serverLog       string
	isServerRunning bool
	runningMutex    sync.Mutex
}

func NewState() *WsState {
	return &WsState{
		ws:              make(map[uint64]*websocket.Conn),
		wsMutex:         sync.Mutex{},
		serverLog:       "",
		isServerRunning: false,
		runningMutex:    sync.Mutex{}}
}

func (wss *WsState) AddConnection(ws *websocket.Conn, id uint64) {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	wss.ws[id] = ws
}

func (wss *WsState) GetConnection(id uint64) *websocket.Conn {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	return wss.ws[id]
}

func (wss *WsState) CloseAndRemoveConnection(id uint64) error {
	wss.wsMutex.Lock()
	defer wss.wsMutex.Unlock()

	err := wss.ws[id].WriteMessage(websocket.CloseNormalClosure, nil)
	delete(wss.ws, id)

	log.Debug().Uint64("request_id", id).Msg("ws closed")

	return err
}

type WsPayload interface {
	Type() string
}

type WsServerState struct {
	Version       string   `json:"version"`
	ConnectUrl    string   `json:"connect_url"`
	ServerName    string   `json:"server_name"`
	IsRunning     bool     `json:"is_running"`
	OnlinePlayers []string `json:"online_players"` // TODO
	AutoStopTime  int64    `json:"auto_stop_time"` // TODO
	BotTag        string   `json:"bot_tag"`        // TODO
}

func (p *WsServerState) Type() string {
	return "state"
}

func NewWsPayloadVersion(conf *config.Config, d *driver.Driver) *WsServerState {
	return &WsServerState{
		Version:       config.Version,
		ConnectUrl:    conf.ConnectUrl,
		ServerName:    d.ServerName(),
		IsRunning:     d.IsRunning(),
		OnlinePlayers: d.OnlinePlayers(),
		AutoStopTime:  time.Now().Add(d.TimeBeforeStop()).Unix(),
		BotTag:        "@my_todo_bot",
	}
}

func serializePayload(payload WsPayload) echo.Map {
	return echo.Map{"type": payload.Type(), "data": payload}
}

func addWebsocket(g *echo.Group, conf *config.Config, driver *driver.Driver) {
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

		err = ws.WriteJSON(serializePayload(NewWsPayloadVersion(conf, driver)))
		if logIfErr(err) {
			return err
		}

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
			case "start":
				err := driver.Start()
				if logIfErr(err) {
					return err
				}
			}
		}
		return nil
	}, authMiddleware)
}

package communication

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
)

const (
	eventData     = "data"
	eventActuator = "actuator"
)

type SocketioServer struct {
	Server       *socketio.Server
	webappSocket socketio.Conn
}

type ServiceDataEvent struct {
	Service string  `json:"service"`
	Data    float64 `json:"data"`
}

type ActuatorConnectedEvent struct {
	IP string `json:"ip"`
}

func NewSocketioServer() (*SocketioServer, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	return &SocketioServer{Server: server}, nil
}

func (ss *SocketioServer) Serve() error {
	ss.Server.OnConnect("/service_data", func(s socketio.Conn) error {
		ss.webappSocket = s
		return nil
	})
	if err := ss.Server.Serve(); err != nil {
		return err
	}
	return nil
}

func (ss *SocketioServer) Close() error {
	if err := ss.Server.Close(); err != nil {
		return err
	}
	return nil
}

func (ss *SocketioServer) ForwardServiceData(evt ServiceDataEvent) {
	jsonData, _ := json.Marshal(&evt)
	ss.webappSocket.Emit(eventData, jsonData)
}

func (ss *SocketioServer) ForwardActuatorData(evt ActuatorConnectedEvent) {
	jsonData, _ := json.Marshal(&evt)
	ss.webappSocket.Emit(eventActuator, jsonData)
}

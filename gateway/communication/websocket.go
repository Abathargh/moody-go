package communication

import (
	"encoding/json"
	"gateway/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	serviceConnection  *websocket.Conn
	actuatorConnection *websocket.Conn

	serviceUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	actuatorUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type SocketioServer struct {
	ServerConn *websocket.Conn
}

type ServiceDataEvent struct {
	Service string  `json:"service"`
	Data    float64 `json:"data"`
}

type ActuatorConnectedEvent struct {
	IP string `json:"ip"`
}

func ServeServiceWS(w http.ResponseWriter, r *http.Request) {
	serviceUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := serviceUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	serviceConnection = conn
}

func ServeActuatorWS(w http.ResponseWriter, r *http.Request) {
	actuatorUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := actuatorUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	actuatorConnection = conn
}

// Marshals and forwards data referring a service event
func ForwardServiceData(evt ServiceDataEvent) {
	jsonData, err := json.Marshal(&evt)
	if err != nil {
		log.Println(err)
		return
	}
	if err := serviceConnection.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println(err)
	}

}

// Gets a mapping of an actuator server that just connected to the gateway and forwards
// it to the webapp
func ForwardActuatorData(evt ActuatorConnectedEvent) {
	// If we don't introduce this delay here, we could be trying to request the mappings
	// before the actuator server is on: the request would fail, and the socketio message
	// would not be sent, but the actuator would eventually turn on: the mappings wouldn't
	// be received without refreshing the webapp.
	time.Sleep(1 * time.Second)
	mappings, err := models.GetActuatorMapping(evt.IP)
	if err != nil {
		log.Println(err)
		return
	}

	actuator := models.Actuator{
		IP:       evt.IP,
		Mappings: mappings.DataTable,
	}

	jsonData, err := json.Marshal(&actuator)
	if err != nil {
		log.Println(err)
		return
	}

	if err := actuatorConnection.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println(err)
	}
}

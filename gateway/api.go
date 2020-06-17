package main

import (
	"gateway/communication"
	"github.com/gorilla/mux"
	"net/http"
)

func HttpListenAndServer(port string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/neural_state", neuralStateMux)
	router.HandleFunc("/actuator_mode", actuatorModeMux)
	router.HandleFunc("/sensor_service", serviceMux)
	router.HandleFunc("/situation", situationMux)
	router.Handle("/socket.io/", communication.SioServer.Server)
	server := &http.Server{Addr: port, Handler: router}
	return server
}

func neuralStateMux(w http.ResponseWriter, r *http.Request) {
	applyHeaders(w)

	switch r.Method {
	case http.MethodPut:
		setNeuralState(w, r)
	case http.MethodOptions:
		respOptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func actuatorModeMux(w http.ResponseWriter, r *http.Request) {
	applyHeaders(w)

	switch r.Method {
	case http.MethodPost:
		setActuatorMode(w, r)
	case http.MethodOptions:
		respOptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func serviceMux(w http.ResponseWriter, r *http.Request) {
	applyHeaders(w)

	switch r.Method {
	case http.MethodPost:
		activateService(w, r)
	case http.MethodGet:
		getActivatedServices(w, r)
	case http.MethodDelete:
		deactivateService(w, r)
	case http.MethodOptions:
		respOptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func situationMux(w http.ResponseWriter, r *http.Request) {
	applyHeaders(w)

	switch r.Method {
	case http.MethodPut:
		setSituation(w, r)
	case http.MethodGet:
		getSituation(w, r)
	case http.MethodOptions:
		respOptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func applyHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Range")
}

func respOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusNoContent)
}

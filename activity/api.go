package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func HttpListenAndServer(port string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/situation", SituationsMux)
	router.HandleFunc("/situation/{id}", SituationMux)
	router.HandleFunc("/service", ServicesMux)
	router.HandleFunc("/service/{id}", ServiceMux)
	router.Use(allowAllCorsMiddleware)
	router.Use(logRequestResponseMiddleWare)
	server := &http.Server{Addr: port, Handler: router}
	return server
}

func allowAllCorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin == "" {
			origin = "*"
		}
		applyHeaders(origin, &w)
		if r.Method == http.MethodOptions {
			respOptions(origin, &w)
		}
		h.ServeHTTP(w, r)
	})
}

func logRequestResponseMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		lrw := NewLoggingResponseWriter(w, r)
		h.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		method := lrw.method
		url := lrw.url

		log.Printf("%s %s %d %s", method, url, statusCode, http.StatusText(statusCode))
	})
}

// Get all the Services or add one
func ServicesMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getServices(w, r)
	case http.MethodPost:
		postService(w, r)
	case http.MethodOptions:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Get/delete a situation or patch it to activate it
func ServiceMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getService(w, r)
	case http.MethodDelete:
		deleteService(w, r)
	case http.MethodOptions:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Get all the situations or add one
func SituationsMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSituations(w, r)
	case http.MethodPost:
		postSituation(w, r)
	case http.MethodOptions:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Get/delete a situation or patch it to activate it
func SituationMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSituation(w, r)
	case http.MethodDelete:
		deleteSituation(w, r)
	case http.MethodOptions:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// CORS Middleware Headers

func applyHeaders(origin string, w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Access-Control-Allow-Headers,Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
	(*w).Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Range")
}

func respOptions(origin string, w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Access-Control-Allow-Headers,Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
	(*w).Header().Set("Access-Control-Max-Age", "1728000")
	(*w).Header().Set("Content-Type", "text/plain; charset=utf-8")
	(*w).Header().Set("Content-Length", "0")
	(*w).WriteHeader(http.StatusNoContent)
}

// Request/Response logger

type loggingResponseWriter struct {
	http.ResponseWriter
	request    *http.Request
	statusCode int
	method     string
	url        string
}

func NewLoggingResponseWriter(w http.ResponseWriter, r *http.Request) *loggingResponseWriter {
	return &loggingResponseWriter{w, r, http.StatusOK, "", ""}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.method = lrw.request.Method
	lrw.url = lrw.request.URL.Path
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

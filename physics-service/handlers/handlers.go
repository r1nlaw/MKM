package handlers

import (
	"net/http"
)

func ForceHandler(w http.ResponseWriter, r *http.Request) {
	forceHandler(w, r)
}

func IntegrateHandler(w http.ResponseWriter, r *http.Request) {
	integrateHandler(w, r)
}

func VectorHandler(w http.ResponseWriter, r *http.Request) {
	vectorHandler(w, r)
}

func TrajectoryHandler(w http.ResponseWriter, r *http.Request) {
	trajectoryHandler(w, r)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	wsHandler(w, r)
}

package handlers

import "net/http"

func IntegrateHandler(w http.ResponseWriter, r *http.Request) {
	integrateHandler(w, r)
}

func CurrentVectorHandler(w http.ResponseWriter, r *http.Request) {
	currentVectorHandler(w, r)
}

func ForceHandler(w http.ResponseWriter, r *http.Request) {
	forceHandler(w, r)
}

func TrajectoryHandler(w http.ResponseWriter, r *http.Request) {
	trajectoryHandler(w, r)
}

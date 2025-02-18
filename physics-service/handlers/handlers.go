package handlers

import "net/http"

func WebSocketConnectionHandle(w http.ResponseWriter, r *http.Request) {
	webSocketConnectionHandler(w, r)
}

func DrowRocketHandler(w http.ResponseWriter, r *http.Request) {
	RocketHandler(w, r)
}

func UpdateRocketThrust(w http.ResponseWriter, r *http.Request) {
	updateRocketThrust(w, r)
}

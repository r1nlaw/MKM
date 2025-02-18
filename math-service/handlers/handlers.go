package handlers

import "net/http"

func IntegrateHandler(w http.ResponseWriter, r *http.Request) {
	integrateHandler(w, r)
}

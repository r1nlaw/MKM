package handlers

import "net/http"

type Integrator interface {
	Integrate(w http.ResponseWriter, r *http.Request)
}

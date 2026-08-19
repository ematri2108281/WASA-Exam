package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// liveness is a simple health check endpoint
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

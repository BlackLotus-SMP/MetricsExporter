package endpoints

import (
	"net/http"
)

type HealthCheck struct {
}

func (h HealthCheck) Route(mux *http.ServeMux) {
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

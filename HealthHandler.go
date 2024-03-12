package routes

import (
	"net/http"
)

// Add here all default behaviour for all/most routes
func HealthHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

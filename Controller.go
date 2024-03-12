package routes

import (
	"net/http"

	"github.com/CodeNamor/Common/server"
	"github.com/gorilla/mux"
)

type Config struct {
	Version     string
	Commit      string
	BuildStamp  string
	BuildNumber string
	Env         string
}

// Calling this from when setting up the router will add these common endpoints. Readiness and liveness update functions returned.
func AddCommonRoutes(router *mux.Router, buildInfo Config, interfacesForExpansion ...interface{}) (readyFn server.UpdateFn, livenessFn server.UpdateFn) {
	router.HandleFunc("/health", HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/version", VersionHandler(buildInfo.Version+buildInfo.BuildStamp,
		buildInfo.Commit, buildInfo.BuildNumber, buildInfo.Env, interfacesForExpansion...)).Methods(http.MethodGet)
	readyFn, livenessFn = server.CreateAndHandleReadinessLiveness(router, "/readiness", "/liveness")
	readyFn(true)
	livenessFn(true)

	return
}

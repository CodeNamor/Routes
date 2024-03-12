package routes

import (
	"testing"

	"github.com/gorilla/mux"
)

func TestAddCommonRoutes(t *testing.T) {

	//TODO: Was considering a test that proves all of the handlers are embedded...
	//router := mux.NewRouter().StrictSlash(true)
	//router.Walk()
}

func routerWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) {
	// Not sure what this is supposed to do... no examples found.
}

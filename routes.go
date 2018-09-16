package gisproxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

//SetRouter  set proxy router
func (server *GISProxy) SetRouter() *GISProxy {
	server.Router = mux.NewRouter()
	return server
}

//SetRoutes set proxy routes
func (server *GISProxy) SetRoutes() *GISProxy {
	server.Router.PathPrefix("/geoserver").Handler(http.StripPrefix("/geoserver", http.HandlerFunc(server.geoserverHandler())))
	return server
}

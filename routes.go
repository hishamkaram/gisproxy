package gisproxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

//SetRouter  set proxy router
func (proxyServer *GISProxy) SetRouter() *GISProxy {
	proxyServer.Router = mux.NewRouter()
	return proxyServer
}

//SetRoutes set proxy routes
func (proxyServer *GISProxy) SetRoutes() *GISProxy {
	proxyServer.Router.HandleFunc("/token-auth", proxyServer.generateTokenHandler).Methods("POST")
	proxyServer.Router.PathPrefix("/geoserver").Handler(http.StripPrefix("/geoserver", http.HandlerFunc(proxyServer.authMiddleware(proxyServer.geoserverHandler))))
	return proxyServer
}

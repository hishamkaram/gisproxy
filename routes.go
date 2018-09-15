package gisproxy

import "net/http"

//SetRoutes set proxy routes
func (server *GISProxy) SetRoutes() *GISProxy {
	server.Router.PathPrefix("/geoserver").Handler(http.StripPrefix("/geoserver", http.HandlerFunc(server.geoserverHandler())))

}

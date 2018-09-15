package gisproxy

import (
	"log"
	"net/http"
)

func (server *GISProxy) geoserverHandler() http.HandlerFunc {
	proxy := server.GetReverseProxy("http://localhost:8080/geoserver")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		r.SetBasicAuth("admin", "geoserver")
		proxy.ServeHTTP(w, r)
	}
}

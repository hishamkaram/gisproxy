package gisproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

//GISProxy  type
type GISProxy struct {
	Geoserver []string
	ArcGIS    []string
	MapServer []string
	Address   string
	SSL       bool
	logger    *logrus.Logger
	DB        string
	Router    *mux.Router
}

//GetReverseProxy return reverse proxy from url
func (server *GISProxy) GetReverseProxy(remoteURL string) *httputil.ReverseProxy {
	remote, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(remote)
}

//SetLogger sets instance logger
func (server *GISProxy) SetLogger() *GISProxy {
	server.logger = GetLogger()
	return server
}

//StartServer GIS Proxy Server
func (server *GISProxy) StartServer() {
	http.Handle("/", server.Router)
	err := http.ListenAndServe(server.Address, nil)
	if err != nil {
		panic(err)
	}
}

package gisproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

//DBConnection  database connection string
type DBConnection struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

//GISProxy  type
type GISProxy struct {
	Geoserver []string
	ArcGIS    []string
	MapServer []string
	Address   string
	SSL       bool
	logger    *logrus.Logger
	DB        *DBConnection
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

//StartProxyServer start new server
func (server *GISProxy) StartProxyServer() {
	server.SetLogger()
	server.SetRouter()
	server.SetRoutes()
	server.MigrateDatabase()
	server.LoadData()
	http.Handle("/", server.Router)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		server.logger.Fatal(err)
	}
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

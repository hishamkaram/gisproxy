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
func (proxyServer *GISProxy) GetReverseProxy(remoteURL string) *httputil.ReverseProxy {
	remote, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(remote)
}

//StartProxyServer start new proxyServer
func (proxyServer *GISProxy) StartProxyServer() {
	proxyServer.SetLogger()
	proxyServer.SetRouter()
	proxyServer.SetRoutes()
	proxyServer.MigrateDatabase()
	proxyServer.LoadData()
	http.Handle("/", proxyServer.Router)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		proxyServer.logger.Fatal(err)
	}
}

//SetLogger sets instance logger
func (proxyServer *GISProxy) SetLogger() *GISProxy {
	proxyServer.logger = GetLogger()
	return proxyServer
}

//StartServer GIS Proxy Server
func (proxyServer *GISProxy) StartServer() {
	http.Handle("/", proxyServer.Router)
	err := http.ListenAndServe(proxyServer.Address, nil)
	if err != nil {
		panic(err)
	}
}

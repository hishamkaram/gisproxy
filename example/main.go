package main

import (
	"github.com/hishamkaram/gisproxy"
)

// func geoserverHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL)
// 		r.SetBasicAuth("admin", "geoserver")
// 		defer r.Body.Close()
// 		p.ServeHTTP(w, r)
// 	}
// }
func main() {
	// remote, err := url.Parse("http://localhost:8080/geoserver")
	// if err != nil {
	// 	panic(err)
	// }
	server := gisproxy.GISProxy{
		DB:      &gisproxy.DBConnection{Name: "gisproxy", Password: "clogic", Username: "hishamkaram", Port: 5432, Host: "localhost"},
		Address: ":8081",
	}
	server.StartProxyServer()

}

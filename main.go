// package gisproxy

// import (
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"

// 	"github.com/gorilla/mux"
// )

// func geoserverHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL)
// 		r.SetBasicAuth("admin", "geoserver")
// 		p.ServeHTTP(w, r)
// 	}
// }
// func main() {
// 	remote, err := url.Parse("http://localhost:8080/geoserver")
// 	if err != nil {
// 		panic(err)
// 	}
// 	proxy := httputil.NewSingleHostReverseProxy(remote)
// 	r := mux.NewRouter()
// 	r.PathPrefix("/geoserver").Handler(http.StripPrefix("/geoserver", http.HandlerFunc(geoserverHandler(proxy))))
// 	http.Handle("/", r)
// 	err = http.ListenAndServe(":8081", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }

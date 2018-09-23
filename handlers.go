package gisproxy

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type jsonError struct {
	Error string `json:"error,omitempty"`
}

//GenerateTokenHandler return token if user is valid
func (proxyServer *GISProxy) generateTokenHandler(w http.ResponseWriter, req *http.Request) {

	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	db, err := proxyServer.GetDB()
	if err != nil {
		proxyServer.logger.Error(err)
	}
	defer db.Close()
	var targetUser User
	if dbc := db.Where(&User{Username: user.Username}).First(&targetUser); dbc.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonError{Error: dbc.Error.Error()})
	} else {
		authenticated := CheckPassword(user.Password, targetUser.Password)
		if authenticated {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": user.Username,
				"password": user.Password,
				"expire":   time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, error := token.SignedString([]byte(proxyServer.SecretKey))
			if error != nil {
				proxyServer.logger.Error(error)
				w.WriteHeader(http.StatusAccepted)
				json.NewEncoder(w).Encode(jsonError{Error: "Internal Server Error"})
			} else {
				w.WriteHeader(http.StatusAccepted)
				json.NewEncoder(w).Encode(AuthToken{Token: tokenString})
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(jsonError{Error: "Authentication Failed"})
		}
	}

}

//userResource return token if user is valid
func (proxyServer *GISProxy) userResource(w http.ResponseWriter, req *http.Request) {
	users, count := proxyServer.getUsers()
	json.NewEncoder(w).Encode(UsersAPIResponse{
		Objects: users, APIResource: APIResource{APIMeta{Count: count}}})
}

//ErrorHandler to generate ogc error
func (proxyServer *GISProxy) errorHandler(w http.ResponseWriter, errMsg string) {
	xml, err := proxyServer.GenerateExceptionReport(&ServiceExceptionReport{Version: "1.1.1", Exceptions: []ServiceException{
		ServiceException{Message: errMsg},
	}})
	if err != nil {
		proxyServer.logger.Error(err)
		xml = []byte(`Internal Server Error 500`)
		http.Error(w, string(xml), http.StatusUnauthorized)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(xml)
	}
}
func (proxyServer *GISProxy) geoserverHandler(w http.ResponseWriter, r *http.Request) {
	r.SetBasicAuth("admin", "geoserver")
	proxy := proxyServer.GetReverseProxy("http://localhost:8080/geoserver")
	proxy.ServeHTTP(w, r)
}

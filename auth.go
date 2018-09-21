package gisproxy

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

//BasicAuth BasicAuth
func (proxyServer *GISProxy) BasicAuth(req *http.Request) (isAuthenticated bool, user *User, authErr error) {
	isAuthenticated = false
	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(auth) == 2 && auth[0] == "Basic" {
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) == 2 {
			db, err := proxyServer.GetDB()
			if err != nil {
				proxyServer.logger.Error(err)
			}
			defer db.Close()
			var targetUser User
			if dbc := db.Where(&User{Username: pair[0]}).First(&targetUser); dbc.Error != nil {
				authErr = dbc.Error
			} else {
				isAuthenticated = CheckPassword(pair[1], targetUser.Password)
				if isAuthenticated {
					user = &targetUser
				} else {
					authErr = errors.New("Authentication Failed")
				}
			}
		}
	}
	return
}

//TokenAuth BasicAuth
func (proxyServer *GISProxy) TokenAuth(req *http.Request) (isAuthenticated bool, user *User, authErr error) {
	isAuthenticated = false
	authorizationHeader := req.Header.Get("authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Token Error")
				}
				return []byte(proxyServer.SecretKey), nil
			})
			if err != nil {
				isAuthenticated = false
				authErr = err
				return
			}
			if token.Valid {
				tokenInfo := token.Claims.(jwt.MapClaims)
				username, _ := tokenInfo["username"]
				db, err := proxyServer.GetDB()
				if err != nil {
					proxyServer.logger.Error(err)
				}
				defer db.Close()
				var targetUser User
				if dbc := db.Where(&User{Username: username.(string)}).First(&targetUser); dbc.Error != nil {
					authErr = dbc.Error
				} else {
					isAuthenticated = true
					user = &targetUser
				}
			} else {
				authErr = errors.New("Invalid Token")
			}
		}
	} else {
		authErr = errors.New("An Authorization Header Is Required")
	}
	return
}

//AuthRequest check if user
func (proxyServer *GISProxy) AuthRequest(req *http.Request) (isAuthenticated bool, user *User, err error) {
	authBackends := []func(*http.Request) (isAuthenticated bool, user *User, authErr error){proxyServer.BasicAuth, proxyServer.TokenAuth}
	for _, backend := range authBackends {
		isAuthenticated, user, err = backend(req)
		if isAuthenticated == true {
			break
		}
	}
	return
}
func (proxyServer *GISProxy) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, _, authErr := proxyServer.AuthRequest(r)
		proxyServer.logger.Info(r.URL)
		if isAuthenticated {
			r.SetBasicAuth("admin", "geoserver")
			proxy := proxyServer.GetReverseProxy("http://localhost:8080/geoserver")
			proxy.ServeHTTP(w, r)
		} else {
			proxyServer.errorHandler(w, authErr.Error())
		}
	}
}

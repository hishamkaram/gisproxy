package gisproxy

import (
	"encoding/base64"
	"net/http"
	"strings"
)

//BasicAuth BasicAuth
func (proxyServer *GISProxy) BasicAuth(req *http.Request) (isAuthenticated bool, user *User) {
	isAuthenticated = false
	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(auth) == 2 && auth[0] == "Basic" {
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		db, err := proxyServer.GetDB()
		if err != nil {
			proxyServer.logger.Error(err)
		}
		defer db.Close()
		if len(pair) == 2 {
			var count uint
			db.Model(&User{}).Where("user_name = ?", pair[0]).Count(&count)
			if count > 0 {
				var targetUser User
				db.Where(&User{UserName: pair[0]}).First(&targetUser)
				isAuthenticated = CheckPassword(pair[1], targetUser.Password)
				if isAuthenticated {
					user = &targetUser
					return
				}
			}
		}
	}
	return
}

//AuthRequest check if user
func (proxyServer *GISProxy) AuthRequest(req *http.Request) (isAuthenticated bool, user *User) {
	authBackends := []func(*http.Request) (isAuthenticated bool, user *User){proxyServer.BasicAuth}
	for _, backend := range authBackends {
		isAuthenticated, user = backend(req)
		if isAuthenticated == true {
			break
		}
	}
	return
}

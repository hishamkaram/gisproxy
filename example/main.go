package main

import (
	"github.com/hishamkaram/gisproxy"
)

func main() {
	server := gisproxy.GISProxy{
		DB:      &gisproxy.DBConnection{Name: "gisproxy", Password: "clogic", Username: "hishamkaram", Port: 5432, Host: "localhost"},
		Address: ":8081",
	}
	server.StartProxyServer()

}

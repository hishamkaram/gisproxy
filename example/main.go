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
	// data, _ := gisproxy.ReadFile("../data/wms/1.1.1/capabilities_1_1_1.xml")
	// cap := gisproxy.ParseWMSCapabilities(data)
	// fmt.Printf("%v", cap)
	// for _, layer := range cap.Capability.Layer.Layer {
	// 	fmt.Printf("\n%v\n", layer)
	// }

}

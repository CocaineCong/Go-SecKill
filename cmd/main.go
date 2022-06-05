package main

import (
	"SecKill/config"
	"SecKill/routers"
)

func main() {
	config.Init()
	r := routers.NewRouter()
	_ = r.Run(config.HttpPort)
}

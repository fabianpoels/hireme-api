package server

import (
	"fmt"
	"hireme-api/config"
)

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port")))
}

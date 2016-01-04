package main

import (
	"flag"

	"photo/server"
	loadConfig "photo/utils/load-config"
	"photo/utils/logs"
)

var (
	flConfigFile = flag.String("config-file", "config-default.json", "Load config from file")

	l = logs.New("photo-server")
)

func main() {
	flag.Parse()

	var cfg server.Config
	err := loadConfig.FromFileAndEnv(&cfg, *flConfigFile)
	if err != nil {
		l.Fatalln("Error loading config:", err)
	}

	server.Start(cfg)
}

package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	db       string
	password string
	host     string
	port     string
	name     string
}

func getConfig() Configuration {
	conf := Configuration{}
	gonfig.GetConf("public/config/config.json", &conf)
	return conf
}

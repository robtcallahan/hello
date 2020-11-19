package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	PlexServer string `json:"plex-server"`
	PlexPort   int    `json:"plex-port"`
	PlexToken  string `json:"plex-token"`
}

func readConfig() *Config {
	contents, err := ioutil.ReadFile("config.json")
	checkError(err)

	var config Config = Config{}
	err = json.Unmarshal([]byte(contents), &config)
	checkError(err)
	return &config
}

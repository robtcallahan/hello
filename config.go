package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	PlexServer      string   `json:"plex-server"`
	PlexPort        int      `json:"plex-port"`
	PlexToken       string   `json:"plex-token"`
	ImageBaseURL    string   `json:"image-base-url"`
	MailFrom        string   `json:"mail-from"`
	MailPassword    string   `json:"mail-password"`
	MailTo          []string `json:"mail-to"`
	RemoteServer    string   `json:"remote-server"`
	RemoteServerKey string   `json:"remote-server-key"`
}

func readConfig() *Config {
	contents, err := ioutil.ReadFile("config.json")
	checkError(err)

	var config Config = Config{}
	err = json.Unmarshal([]byte(contents), &config)
	checkError(err)
	return &config
}

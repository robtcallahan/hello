package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Location ...
type Location struct {
	ID   int    `xml:"id,attr"`
	Path string `xml:"path,attr"`
}

// Directory ...
type Directory struct {
	Key      int      `xml:"key,attr"`
	Type     string   `xml:"type,attr"`
	Title    string   `xml:"title,attr"`
	Location Location `xml:"Location"`
}

// Genre ...
type Genre struct {
	Tag string `xml:"tag,attr"`
}

// Writer ...
type Writer struct {
	Tag string `xml:"tag,attr"`
}

// Director ...
type Director struct {
	Tag string `xml:"tag,attr"`
}

// Part ...
type Part struct {
	ID        int    `xml:"id,attr"`
	Key       string `xml:"key,attr"`
	Duration  int    `xml:"duration,attr"`
	File      string `xml:"file,attr"`
	Size      int    `xml:"size,attr"`
	Container string `xml:"container,attr"`
}

// Media ...
type Media struct {
	ID              int    `xml:"id,attr"`
	VideoResolution string `xml:"videoResolution,attr"`
	Duration        int    `xml:"duration,attr"`
	Container       string `xml:"container,attr"`
	Part            []Part `xml:"part"`
}

// Video ...
type Video struct {
	Key                   string     `xml:"key,attr"`
	Type                  string     `xml:"type,attr"`
	Studio                string     `xml:"studio,attr"`
	Title                 string     `xml:"title,attr"`
	TitleSort             string     `xml:"titleSort,attr"`
	ContentRating         string     `xml:"contentRating,attr"`
	Summary               string     `xml:"summary,attr"`
	Rating                float32    `xml:"rating,attr"`
	Year                  string     `xml:"year,attr"`
	Tagline               string     `xml:"tagline,attr"`
	Duration              int        `xml:"duration,attr"`
	OriginallyAvailableAt string     `xml:"originallyAvailableAt,attr"`
	AddedAt               int        `xml:"addedAt,attr"`
	UpdatedAt             int        `xml:"updatedAt,attr"`
	Media                 Media      `xml:"Media"`
	Genres                []Genre    `xml:"Genre"`
	Writers               []Writer   `xml:"Writer"`
	Directors             []Director `xml:"Director"`
}

// MediaContainer ...
type MediaContainer struct {
	// XMLName   xml.Name  `xml:"MediaContainer"`
	Size             int         `xml:"size,attr"`
	Title            string      `xml:"title1,attr"`
	LibrarySectionID int         `xml:"librarySectionID"`
	Directories      []Directory `xml:"Directory"`
	Videos           []Video     `xml:"Video"`
}

// PlexClient ...
type PlexClient struct {
	Server     string
	Port       int
	BaseURL    string
	XPlexToken string
	MoviesKey  int
	Client     *http.Client
}

// NewPlexClient ...
func NewPlexClient(c *Config) *PlexClient {
	return &PlexClient{
		XPlexToken: c.PlexToken,
		Server:     c.PlexServer,
		Port:       c.PlexPort,
		BaseURL:    fmt.Sprintf("http://%s:%d", c.PlexServer, c.PlexPort),
		Client:     &http.Client{},
	}
}

// GetLibraries ...
func (c *PlexClient) GetLibraries() *MediaContainer {
	path := fmt.Sprintf("%s/library/sections/?X-Plex-Token=%s", c.BaseURL, c.XPlexToken)
	resp := c.get(path)

	var media MediaContainer
	body, _ := ioutil.ReadAll(resp.Body)
	err := xml.Unmarshal(body, &media)
	checkError(err)

	for _, d := range media.Directories {
		if d.Type == "movie" && d.Title == "Movies" {
			c.MoviesKey = d.Key
		}
	}
	return &media
}

// GetMovies ...
func (c *PlexClient) GetMovies() *MediaContainer {
	path := fmt.Sprintf("%s/library/sections/%d/all?X-Plex-Token=%s", c.BaseURL, c.MoviesKey, c.XPlexToken)
	resp := c.get(path)

	var media MediaContainer
	body, _ := ioutil.ReadAll(resp.Body)
	err := xml.Unmarshal(body, &media)
	checkError(err)
	return &media
}

func (c *PlexClient) get(path string) *http.Response {
	req, err := http.NewRequest("GET", path, nil)
	checkError(err)

	resp, err := c.Client.Do(req)
	checkError(err)
	return resp
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

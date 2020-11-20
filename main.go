package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"sort"
	"time"
)

func main() {
	config := readConfig()
	client := NewPlexClient(config)

	client.GetLibraries()
	media := client.GetMovies()
	media.ServerConfig = ServerConfig{
		PlexToken:       config.PlexToken,
		ImageBaseURL:    config.ImageBaseURL,
		RemoteServer:    config.RemoteServer,
		RemoteServerKey: config.RemoteServerKey,
	}

	sort.Slice(media.Videos, func(i, j int) bool {
		return media.Videos[i].AddedAt < media.Videos[j].AddedAt
	})

	html := outputVideos(media)
	ioutil.WriteFile("videos.html", html.Bytes(), 0644)
	sendMail(config, html)
}

func outputVideos(media *MediaContainer) bytes.Buffer {
	var buf bytes.Buffer
	vids := []Video{}
	minus7days := minusDays(90)

	for _, v := range media.Videos {
		if v.AddedAt > minus7days {
			v.DurationFormatted = fmtDuration(v.Duration)
			vids = append(vids, v)
		}
	}

	sort.Slice(vids, func(i, j int) bool {
		return vids[i].Title < vids[j].Title
	})

	media.Videos = vids

	tmpl := template.Must(template.ParseFiles("videos_tmpl.html")).Option("missingkey=error")
	tmpl.Execute(&buf, media)
	return buf
}

func sendMail(config *Config, html bytes.Buffer) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: New Plex Movies Available!\n"
	msg := []byte(subject + mime + html.String())

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", config.MailFrom, config.MailPassword, "smtp.gmail.com"),
		config.MailFrom,
		config.MailTo,
		[]byte(msg))

	if err != nil {
		panic(err)
	}
}

func minus1day() int64 {
	return time.Now().Unix() - int64(60*60*24)
}

func minusDays(days int64) int64 {
	return time.Now().Unix() - int64(60*60*24*days)
}

func secs2Date(secs int64) string {
	t := time.Unix(secs, 0)
	return t.Format("2006-01-02 15:04:05")
}

func fmtDuration(ms int) string {
	modTime := time.Now().Round(0).Add(-time.Duration(ms) * time.Millisecond)
	d := time.Since(modTime)
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

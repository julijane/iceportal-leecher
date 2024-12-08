package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type AudiobookEpisode struct {
	SerialNumber int    `json:"serialNumber"`
	Title        string `json:"title"`
	Path         string `json:"path"`
}

type Audiobook struct {
	Title string             `json:"title"`
	Files []AudiobookEpisode `json:"files"`
}

type File struct {
	Path string `json:"path"`
}

func (l *Leecher) fetchAllAudiobooks() {
	var audiobooks TeaserList
	err := l.getJson("https://iceportal.de/api1/rs/page/hoerbuecher", &audiobooks)
	if err != nil {
		log.Fatal(err)
	}
	for _, audiobook := range audiobooks.TeaserGroups[0].Items {
		audiobookID := strings.Split(audiobook.Navigation.Href, "/")[2]
		log.Print("Fetching ", audiobook.Title, " (", audiobookID, ")")
		l.fetchAudiobook(audiobookID)
	}
}

func (l *Leecher) fetchAudiobook(audiobookID string) {
	var audiobook Audiobook
	err := l.getJson("https://iceportal.de/api1/rs/page/hoerbuecher/"+audiobookID, &audiobook)
	if err != nil {
		log.Fatal(err)
	}

	dirPath := path.Join(
		"audiobooks",
		sanitizeFileOrPathName(audiobook.Title),
	)

	if err = os.MkdirAll(dirPath, 0o0755); err != nil {
		log.Fatal("Creating directory for audiobook: ", err)
	}

	for _, episode := range audiobook.Files {
		log.Println("Fetching episode", episode.SerialNumber, "-", episode.Title)

		var file File
		err := l.getJson("https://iceportal.de/api1/rs/audiobooks/path"+episode.Path, &file)
		if err != nil {
			log.Fatal(err)
		}

		episodeFilename := fmt.Sprintf("%03d", episode.SerialNumber) + "_" +
			sanitizeFileOrPathName(episode.Title) + ".mp4"

		err = l.saveTo(
			"https://iceportal.de"+file.Path,
			path.Join(
				dirPath,
				episodeFilename,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

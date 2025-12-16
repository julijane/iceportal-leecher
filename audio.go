package main

import (
	"fmt"
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

func (l *Leecher) fetchAllAudio(subsection string) {
	var audiolist TeaserList
	err := l.getJson("https://iceportal.de/api1/rs/page/"+subsection, &audiolist)
	if err != nil {
		sugar.Fatal(err)
	}
	for _, audio := range audiolist.TeaserGroups[0].Items {
		audioID := strings.Split(audio.Navigation.Href, "/")[2]
		sugar.Infof("Fetching %s (%s)", audio.Title, audioID)
		l.fetchAudio(subsection, audioID)
	}
}

func (l *Leecher) fetchAudio(subsection string, audiobookID string) {
	var audiobook Audiobook
	err := l.getJson("https://iceportal.de/api1/rs/page/"+subsection+"/"+audiobookID, &audiobook)
	if err != nil {
		sugar.Errorf("Fetching audiobook metadata: %v", err)
		return
	}

	dirPath := path.Join(
		subsection,
		sanitizeFileOrPathName(audiobook.Title),
	)

	if err = os.MkdirAll(dirPath, 0o0755); err != nil {
		sugar.Fatalf("Creating directory for audiobook: %v", err)
	}

	for _, episode := range audiobook.Files {
		sugar.Infof("Fetching episode %d - %s", episode.SerialNumber, episode.Title)

		var file File
		err := l.getJson("https://iceportal.de/api1/rs/"+episode.Path, &file)
		if err != nil {
			sugar.Fatal(err)
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
			sugar.Fatal(err)
		}
	}
}

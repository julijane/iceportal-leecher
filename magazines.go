package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Magazine struct {
	Title      string `json:"title"`
	Paymethod  string `json:"paymethod"`
	Date       string `json:"date"`
	Navigation struct {
		Href string `json:"href"`
	} `json:"navigation"`
}

func (l *Leecher) fetchAllMagazines() {
	var magazines TeaserList
	err := l.getJson("https://iceportal.de/api1/rs/page/zeitungskiosk", &magazines)
	if err != nil {
		panic(err)
	}

	for _, magazine := range magazines.TeaserGroups[0].Items {
		magazineID := strings.Split(magazine.Navigation.Href, "/")[2]
		log.Print("Fetching ", magazine.Title, " (", magazineID, ")")
		l.fetchMagazine(magazineID)
	}
}

func (l *Leecher) fetchMagazine(magazineID string) {
	var magazine Magazine
	err := l.getJson("https://iceportal.de/api1/rs/page/zeitungskiosk/"+magazineID, &magazine)
	if err != nil {
		panic(err)
	}

	if magazine.Paymethod != "free" && magazine.Paymethod != "freeCopy" {
		log.Fatal("Not free, skipping. Paymethod: ", magazine.Paymethod)
	}

	dirPath := path.Join("magazines", sanitizeFileOrPathName(magazine.Title))

	resp, err := http.Get("https://iceportal.de/" + magazine.Navigation.Href)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(dirPath, 0o0755); err != nil {
		log.Fatal("Creating directory for magazine: ", err)
	}

	err = l.saveTo(
		"https://iceportal.de/"+magazine.Navigation.Href,
		path.Join(
			dirPath,
			sanitizeFileOrPathName(magazine.Title)+"_"+magazine.Date+".pdf",
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}

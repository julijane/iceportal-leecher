package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/k0kubun/pp/v3"
)

type Leecher struct {
	cookie string
}

type TeaserList struct {
	TeaserGroups []struct {
		Items []struct {
			Title      string `json:"title"`
			Navigation struct {
				Href string `json:"href"`
			} `json:"navigation"`
		} `json:"items"`
	} `json:"teaserGroups"`
}

func main() {
	_ = pp.Print

	l := Leecher{}
	l.fetchCookie()
	l.fetchAllAudiobooks()
	l.fetchAllMagazines()
}

func (l *Leecher) fetchCookie() {
	// make a request to the start page to obtain the session tookie cooen
	resp, err := http.Get("https://iceportal.de")
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	cookieHeader := resp.Header.Get("Set-Cookie")
	if cookieHeader == "" {
		log.Fatal("No cookie header found")
	}

	l.cookie = strings.Split(cookieHeader, ";")[0]
}

func (l *Leecher) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if l.cookie != "" {
		req.Header.Set("Cookie", l.cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func (l *Leecher) getJson(url string, v any) error {
	resp, err := l.get(url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	if err != nil {
		return err
	}

	return nil
}

func (l *Leecher) saveTo(url string, filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		log.Print("File ", filePath, " exists already, skipping")
		return nil
	}

	resp, err := l.get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	io.Copy(outFile, resp.Body)

	log.Print("Saved to ", filePath)
	return nil
}

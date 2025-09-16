package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ItemParse struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Genres  string `json:"genres"`
	Img     string `json:"img"`
	Types   string `json:"types"`
	Episode string `json:"episode"`
}

var Items []ItemParse

func GetItem(i int, s *goquery.Selection) {
	// For each item found, get the title
	a := s.Find("a.film-item")
	href, _ := a.Attr("href")
	title := s.Find(".film-item-title").Text()
	genres := s.Find(".film-item-genres").Text()

	film := s.Find(".film-item-image")
	img, _ := film.Find("img").Attr("data-src")
	types := film.Find(".film-item-type").Text()
	episode := film.Find(".film-item-episode").Text()

	item := ItemParse{
		Title:   title,
		Href:    href,
		Genres:  genres,
		Img:     img,
		Types:   types,
		Episode: episode,
	}

	Items = append(Items, item)

	p(3, " → ", "[+]", i, href, title, genres)
	p(2, " → ", "[+]", strings.Repeat("~", 39))
}

func GetHtml(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		e := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(e)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetScrape(url string) []ItemParse {
	// Request the HTML page.
	doc, err := GetHtml(url)
	if err != nil {
		log.Fatal(err)
	}

	Items = nil

	// Find the review items
	doc.Find("#pdopage > div.catalog-list > div.catalog-list-item").Each(GetItem)

	return Items
}

func main() {
	Url := "https://jetfilm.net/serials/"
	FilmItems := GetScrape(Url)

	// add json
	writeJson(FilmItems, "./json/item.json")
}

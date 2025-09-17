package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ItemParse struct {
	Title   string   `json:"title"`
	Href    string   `json:"href"`
	Genres  []string `json:"genres"`
	Img     string   `json:"img"`
	Types   string   `json:"types"`
	Episode Episode  `json:"episode"`
	Rating  Rating   `json:"rating"`
}

type Rating struct {
	Kp   float64 `json:"kp"`
	Imdb float64 `json:"imdb"`
	Jet  float64 `json:"jet"`
}

type Episode struct {
	Season int `json:"season"`
	Series int `json:"series"`
}

var Items []ItemParse

func SetFloat(s string) float64 {
	if s == "" {
		return 0
	}

	s = strings.TrimSpace(s)
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return num
}

func GetRating(film *goquery.Selection) Rating {
	rating := film.Find(".film-item-rating")
	kp := rating.Find(".kp").Text()
	imdb := rating.Find(".imdb").Text()
	jet := rating.Find(".jet").Text()

	return Rating{Kp: SetFloat(kp), Imdb: SetFloat(imdb), Jet: SetFloat(jet)}
}

func GetEpisode(film *goquery.Selection) (Episode, error) {
	episode := film.Find(".film-item-episode").Text()
	//1 сезон 10 серия
	epi := strings.TrimSpace(episode)
	if epi == "" {
		return Episode{}, errors.New("not episode")
	}

	var Season int
	var Series int
	episodeArr := strings.Split(epi, " ")
	if len(episodeArr[0]) != 0 {
		season, err := strconv.Atoi(episodeArr[0])
		if err != nil {
			return Episode{}, errors.New("not season")
		}

		Season = season
	}

	if len(episodeArr[2]) != 0 {
		series, err := strconv.Atoi(episodeArr[2])
		if err != nil {
			return Episode{}, errors.New("not series")
		}

		Series = series
	}

	return Episode{Season: Season, Series: Series}, nil
}

func GetItem(i int, s *goquery.Selection) {
	// For each item found, get the title
	a := s.Find("a.film-item")
	href, _ := a.Attr("href")
	title := s.Find(".film-item-title").Text()
	genres := s.Find(".film-item-genres").Text()

	film := s.Find(".film-item-image")
	img, _ := film.Find("img").Attr("data-src")
	types := film.Find(".film-item-type").Text()

	item := ItemParse{
		Title:  title,
		Href:   href,
		Img:    img,
		Types:  strings.TrimSpace(types),
		Rating: GetRating(film),
	}

	episode, err := GetEpisode(film)
	if err == nil {
		item.Episode = episode
	}

	if len(genres) != 0 {
		item.Genres = strings.Split(genres, ", ")
	}

	Items = append(Items, item)

	p(2, " → ", "[+]", i, href, title, genres)
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
	UrlPage := Url
	f := "./json/item.json"
	for v := range 3 {
		if v > 0 {
			UrlPage = fmt.Sprintf("%s?page=%d", Url, v+1)
		}

		FilmItems := GetScrape(UrlPage)

		p(3, " ~ ", "[+]", v+1, UrlPage, len(FilmItems))

		// add json
		dataFiles, err := loadJson(f)
		if err != nil {
			writeJson(FilmItems, f)
		} else if len(dataFiles) > 0 {
			dataFiles = append(dataFiles, FilmItems...)
			writeJson(dataFiles, f)
		}
	}
}

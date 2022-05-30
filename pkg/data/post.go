package data

import (
	"fmt"
	goose "github.com/advancedlogic/GoOse"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Entry struct {
	Title  string
	URL    string
	Year   int
	Author string
}

var postReader = goose.New()

func FromLinkList(path string) (Set, error) {
	var input Set

	entries, err := getEntries(path)
	if err != nil {
		return input, err
	}
	input.Mapping = Mapping{
		Path: path,
		Fields: map[string]string{
			"title":  "string",
			"text":   "string",
			"author": "string",
			"date":   "numeric",
			"tags":   "keywords",
		},
	}

	sources, err := readEntries(entries)
	if err != nil {
		return input, err
	}
	input.Sources = sources

	return input, nil
}

func getEntries(key string) ([]Entry, error) {
	var entries []Entry

	src, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(src, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func readEntries(entries []Entry) ([]Source, error) {
	var sources []Source

	for _, entry := range entries {
		article, err := postReader.ExtractFromURL(entry.URL)
		if err != nil {
			fmt.Println("Failed", entry.Title)
			return sources, err
		}

		year, _, _ := time.Now().Date()
		if article.PublishDate != nil {
			year, _, _ = article.PublishDate.Date()
		}

		_id := fmt.Sprintf(
			"title=%s&url=%s&author=%s&year=%s",
			article.Title, article.FinalURL, entry.Author, year)

		sources = append(sources, Source{
			ID:  _id,
			URL: article.FinalURL,
			Fields: map[string]interface{}{
				"title":  article.Title,
				"text":   article.CleanedText,
				"author": entry.Author,
				"year":   year,
				"tags":   article.Tags,
			},
		})
	}

	return sources, nil
}

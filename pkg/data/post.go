package data

import (
	"encoding/json"
	"fmt"
	goose "github.com/advancedlogic/GoOse"
	"io/ioutil"
)

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
			"type":   "string",
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

	err = json.Unmarshal(src, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func readEntries(entries []Entry) ([]Source, error) {
	var sources []Source

	for _, entry := range entries {
		body := entry.Title
		if entry.Type == "post" {
			article, err := postReader.ExtractFromURL(entry.URL)
			if err != nil {
				fmt.Println("Failed", entry.Title)
				return sources, err
			}
			body = article.CleanedText
		}

		_id := fmt.Sprintf(
			"title=%s&url=%s&author=%s&year=%d",
			entry.Title, entry.URL, entry.Author, entry.Year)

		sources = append(sources, Source{
			ID:  _id,
			URL: entry.URL,
			Fields: map[string]interface{}{
				"title":  entry.Title,
				"text":   body,
				"author": entry.Author,
				"year":   entry.Year,
				"type":   entry.Type,
			},
		})
	}

	return sources, nil
}

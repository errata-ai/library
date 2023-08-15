package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/k3a/html2text"
)

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

	src, err := os.ReadFile(key)
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
			response, err := http.Get(entry.URL)
			if err != nil {
				fmt.Println("Failed", entry.Title, err)
				return sources, err
			}
			defer response.Body.Close()

			bytes, err := io.ReadAll(response.Body)
			if err != nil {
				return sources, err
			}
			body = html2text.HTML2Text(string(bytes))
		}

		_id := fmt.Sprintf(
			"title=%s&url=%s&author=%s&year=%d&type=%s",
			entry.Title, entry.URL, entry.Author, entry.Year, entry.Type)

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

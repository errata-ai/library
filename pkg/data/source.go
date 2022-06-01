package data

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"gopkg.in/yaml.v2"
)

type Entry struct {
	Title  string
	URL    string
	Year   int
	Author string
	Type   string
	Body   string
}

type Mapping struct {
	Path      string
	Template  string            // TextFSM, LOIS, etc.
	Fields    map[string]string // Key -> type
	Extension string
}

// Source represents an arbitrary, indexable source of data: it could be
// a structured file (YAML, JSON, CSV, XLSX, etc.) or a free-form, plain-text
// file.
type Source struct {
	ID     string
	URL    string
	Fields interface{}
}

// Set represents a data set of arbitrary structure.
type Set struct {
	Sources []Source
	Mapping Mapping
}

// GetDocMapping creates a custom mapping for our DataSource structure.
func GetDocMapping(keys map[string]string) (*mapping.DocumentMapping, error) {
	// TODO: Choose mapping ...
	//
	// See http://bleveanalysis.couchbase.com/analysis
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Store = false
	keywordFieldMapping.Analyzer = keyword.Name

	simpleFieldMapping := bleve.NewTextFieldMapping()
	simpleFieldMapping.Store = false
	simpleFieldMapping.Analyzer = simple.Name

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Store = true
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	stdTextFieldMapping := bleve.NewTextFieldMapping()
	stdTextFieldMapping.Store = false
	stdTextFieldMapping.Analyzer = standard.Name

	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	dateFieldMapping.Store = false

	numFieldMapping := bleve.NewNumericFieldMapping()
	numFieldMapping.Store = false

	custom := bleve.NewDocumentMapping()

	for k, v := range keys {
		if strings.TrimSpace(k) != "" {
			switch v {
			case "datetime":
				custom.AddFieldMappingsAt(k, dateFieldMapping)
			case "numeric":
				custom.AddFieldMappingsAt(k, numFieldMapping)
			case "keywords":
				custom.AddFieldMappingsAt(k, keywordFieldMapping)
			default:
				// NOTE: Should we go with "simple" in hopes of keeping
				// lang-agnostic ...
				//
				// en is smaller/faster ...
				custom.AddFieldMappingsAt(k, englishTextFieldMapping)
			}
		}
	}

	return custom, nil
}

func ParseMap(path string) (Mapping, error) {
	var mapped Mapping

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return mapped, err
	}

	err = yaml.Unmarshal(src, &mapped)
	if err != nil {
		return mapped, err
	}

	return mapped, nil
}

func FromMapping(path string) (Set, error) {
	var input Set

	m, err := ParseMap(path)
	if err != nil {
		return input, err
	}
	input.Mapping = m

	sources, err := FromFile(input.Mapping)
	if err != nil {
		return input, err
	}
	input.Sources = sources

	return input, nil
}

func FromFile(m Mapping) ([]Source, error) {
	var input []Source

	if ret, err := isDir(m.Path); ret && err == nil {
		return fromFolder(m)
	}

	src, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return input, err
	}

	switch filepath.Ext(m.Path) {
	case ".tsv":
		return fromTSV(src, m.Fields)
	case ".csv":
		return fromCSV(src, m.Fields)
	}

	return input, nil
}

func Context(id, path string) (string, error) {
	if ret, err := isDir(path); ret && err == nil {
		return extractFile(path, id)
	}

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	switch filepath.Ext(path) {
	case ".tsv":
		return extractTSV(src, id)
	case ".csv":
		return extractCSV(src, id)
	}

	return "", nil
}

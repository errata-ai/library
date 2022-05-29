package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/errata-ai/library/pkg/data"
)

// Engine implements our local search engine interface.
type Engine struct {
	Index bleve.Index
	Count uint64
}

// Match represents an individual search match.
//
// It consists of an ID (row number, file path, etc.) and, optionally, the
// in-text context.
type Match struct {
	ID        string
	Locations search.FieldTermLocationMap
	Score     string
	Tokens    string
	Fragments []string
}

// Result is the product of performing a query.
type Result struct {
	Fragments search.FieldFragmentMap
	URL       string
	ID        string
}

// BulkLoad performs batch indexing.
func (e *Engine) BulkLoad(src []data.Source) error {
	var err error

	batch := e.Index.NewBatch()
	bsize := 0

	for _, r := range src {
		if r.ID != "" {
			err = batch.Index(r.ID, r.Fields)
			if err != nil {
				return err
			}
			bsize++
			// TODO: Batch sizing ...
			//
			// all: 43.683 total
			// 100: 1m32s
			// 1k: 36.366
			// 2k: 35.2
			if bsize >= 1000 {
				err = e.Index.Batch(batch)
				if err != nil {
					return err
				}
				batch = e.Index.NewBatch()
				bsize = 0
			}
		}
	}

	if bsize > 0 {
		err = e.Index.Batch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

// Search the index for the given query.
func (e *Engine) Search(q string) (search.DocumentMatchCollection, error) {
	// The `q` is the "human-friendly" query syntax:
	//
	// +foo (required have)
	// -foo (require not have)
	// foo~{1,2} (fuzzy)
	// "foo bar" (phrase)
	// foo bar (OR)
	//
	// ... also possible to have field-specific and duration comparisons.
	seq := bleve.NewQueryStringQuery(q)

	req := bleve.NewSearchRequestOptions(seq, int(e.Count), 0, false)
	req.Highlight = bleve.NewHighlightWithStyle("html")

	results, err := e.Index.Search(req)
	if err != nil {
		return nil, err
	}

	return results.Hits, nil
}

// NewEngine creates a new search engine with the given name.
func NewEngine(name string, keys map[string]string) (*Engine, error) {
	var engine Engine

	idxMapping := bleve.NewIndexMapping()

	mapping, err := data.GetDocMapping(keys)
	if err != nil {
		return &engine, err
	}
	idxMapping.DefaultMapping = mapping

	index, err := bleve.New(name, idxMapping)
	if err != nil {
		return &engine, err
	}

	engine.Index = index
	count, err := index.DocCount()
	if err != nil {
		return &engine, err
	}
	engine.Count = count

	return &engine, nil
}

// NewEngineFromData creates a new search engine with the given name and data.
func NewEngineFromData(name string, data data.Set) (*Engine, error) {
	engine, err := NewEngine(name, data.Mapping.Fields)
	if err != nil {
		return engine, err
	}

	err = engine.BulkLoad(data.Sources)
	if err != nil {
		return engine, err
	}

	return engine, nil
}

// LoadEngine loads an engine from disk.
func LoadEngine(name string) (*Engine, error) {
	index, err := bleve.Open(name)
	if err != nil {
		return &Engine{}, err
	}

	count, err := index.DocCount()
	if err != nil {
		return &Engine{}, err
	}

	return &Engine{Index: index, Count: count}, nil
}

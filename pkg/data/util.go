package data

import (
	"bytes"
	"encoding/csv"
	"errors"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/errata-ai/library/internal/nlp"
)

func hash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func toNum(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func toValue(s, k string) (interface{}, error) {
	switch k {
	case "numeric":
		return toNum(s)
	case "datetime":
		// TODO: This appears to be slow?
		//
		// should we require date format upfront?
		return dateparse.ParseAny(s)
	default:
		return s, nil
	}
}

func fromFolder(m Mapping) ([]Source, error) {
	var input []Source

	files, err := ioutil.ReadDir(m.Path)
	if err != nil {
		return input, err
	}

	for idx, file := range files {
		if filepath.Ext(file.Name()) == m.Extension {
			p := filepath.Join(m.Path, file.Name())
			if err != nil {
				return input, err
			}

			if m.Template != "" {
				// We have a template ...
				out, err := nlp.DoSectionize(p, m.Template)
				// ^ THIS is the new file -- essentially, a JSON input.
				if err != nil {
					return input, err
				}

				for _, section := range out.Dict {
					input = append(input, Source{
						ID:     strconv.Itoa(idx + 1),
						Fields: section,
					})
				}
			}
		}
	}

	return input, nil
}

func fromSeparated(src []byte, sep rune, m map[string]string) ([]Source, error) {
	var input []Source
	var headr []string

	r := csv.NewReader(bytes.NewReader(src))
	r.Comma = sep

	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return input, err
		} else if count == 0 {
			headr = record
		} else if count > 0 {
			entry := map[string]interface{}{}
			for idx, k := range headr {
				v, err := toValue(record[idx], m[k])
				if err != nil {
					return input, err
				}
				entry[k] = v
			}
			input = append(input, Source{
				ID:     strconv.Itoa(count + 1),
				Fields: entry,
			})
		}
		count++
	}

	return input, nil
}

func fromTSV(src []byte, m map[string]string) ([]Source, error) {
	return fromSeparated(src, '\t', m)
}

func fromCSV(src []byte, m map[string]string) ([]Source, error) {
	return fromSeparated(src, ',', m)
}

func extractSeparated(src []byte, id string, sep rune) (string, error) {
	r := csv.NewReader(bytes.NewReader(src))
	r.Comma = sep

	row, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		} else if count == row-1 {
			return strings.Join(record, "\n\n"), nil
		}
		count++
	}

	return "", errors.New("not found")
}

func extractFile(path, id string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for idx, file := range files {
		p := filepath.Join(path, file.Name())
		if strconv.Itoa(idx+1) == id {
			src, err := ioutil.ReadFile(p)
			if err != nil {
				return "", err
			}
			return string(src), nil
		}
	}

	return "", errors.New("not found")
}

func extractTSV(src []byte, id string) (string, error) {
	return extractSeparated(src, id, '\t')
}

func extractCSV(src []byte, id string) (string, error) {
	return extractSeparated(src, id, ',')
}

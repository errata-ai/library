// See also
//
// http://textfsm.nornir.tech/
// https://github.com/dmulyalin/ttp
// https://github.com/google/textfsm
// https://github.com/sirikothe/gotextfsm
package line

import (
	"strings"
)

type FieldMap = map[string][]string

type Key struct {
	Type string
	Name string
}

func ParseLines(text string, keys []Key) (FieldMap, error) {
	fields := FieldMap{}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	jdx := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		key := keys[jdx]

		fields[key.Name] = append(fields[key.Name], line)
		if key.Type == "single" {
			jdx++
		}
	}

	if jdx < len(keys) {
		tooMany := keys[jdx]
		pool := fields[tooMany.Name]

		pos := len(keys)
		m := false

		count := 0
		limit := len(keys) - jdx

		for pos >= limit {
			if !m || len(pool) == limit+1 {
				pos--
			}

			last := pool[len(pool)-1]
			if len(pool) > 0 {
				pool = pool[:len(pool)-1]
			}
			fields[tooMany.Name] = pool

			key := keys[pos]
			fields[key.Name] = append([]string{last}, fields[key.Name]...)
			if key.Type == "multiple" {
				m = true
			} else {
				m = false
			}

			count++
		}
	}

	return fields, nil
}

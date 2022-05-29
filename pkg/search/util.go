package search

import (
	"encoding/json"
	"math"
)

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func JSONEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b[1 : len(b)-1])
}

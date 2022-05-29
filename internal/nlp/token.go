package nlp

import (
	"strconv"
	"strings"

	"github.com/GolangNLP/split"
)

type TokenInfo struct {
	Text string
	Type string
	Info string
}

func TextToTokens(s string) ([]TokenInfo, error) {
	var tokens []TokenInfo

	tok, err := split.NewIterTokenizer()
	if err != nil {
		return tokens, err
	}

	for _, t := range tok.Split(s) {
		if strings.TrimSpace(t) == "" {
			continue
		} else if _, ok := EnglishStopWords[strings.ToLower(t)]; ok {
			tokens = append(tokens, TokenInfo{
				Text: t,
				Type: "STOP",
				Info: "Stop Word (excluded from search)",
			})
		} else if _, err := strconv.ParseFloat(t, 64); err == nil {
			tokens = append(tokens, TokenInfo{
				Text: t,
				Type: "NUM",
				Info: "Number",
			})
		} else if !IsLetter(t) {
			tokens = append(tokens, TokenInfo{
				Text: t,
				Type: "SYM",
				Info: "Symbol",
			})
		} else {
			tokens = append(tokens, TokenInfo{
				Text: t,
				Type: "WORD",
				Info: "Word",
			})
		}
	}

	return tokens, nil
}

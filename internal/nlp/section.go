package nlp

import (
	"io/ioutil"

	"github.com/sirikothe/gotextfsm"
)

func DoSectionize(s, t string) (gotextfsm.ParserOutput, error) {
	var parsed gotextfsm.ParserOutput

	input, err := ioutil.ReadFile(s)
	if err != nil {
		return parsed, err
	}

	temp, err := ioutil.ReadFile(t)
	if err != nil {
		return parsed, err
	}

	return Sectionize(string(input), string(temp))
}

// Sectionize splits semi-structured text into sections according to the
// user-defined
func Sectionize(s, t string) (gotextfsm.ParserOutput, error) {
	var parsed gotextfsm.ParserOutput
	var fsm gotextfsm.TextFSM

	err := fsm.ParseString(t)
	if err != nil {
		return parsed, err
	}

	err = parsed.ParseTextString(s, fsm, true)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

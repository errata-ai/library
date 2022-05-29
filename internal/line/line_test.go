package line

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

var testdata = filepath.Join("..", "..", "testdata")

func TestParseEmail(t *testing.T) {
	text, err := ioutil.ReadFile(filepath.Join(testdata, "email.txt"))
	if err != nil {
		t.Error(err)
	}

	fields, err := ParseLines(string(text), []Key{
		{Type: "single", Name: "opening"},
		{Type: "multiple", Name: "body"},
		{Type: "single", Name: "closing"},
		{Type: "multiple", Name: "signature"},
	})
	if err != nil {
		t.Error(err)
	}

	expected := map[string][]string{
		"body": {
			"I'm interesting in a custom oil painting of a dog (similar to https://shop.oldiart.com/en/spd/01em/Egyedi-festmenyrendeles), but I had two questions:",
			"1. Do you take orders from the U.S.?",
			"2. How much does it cost?",
		},
		"closing": {"Thank you,"},
		"opening": {"Hi,"},
		"signature": {
			"Joseph",
			"B.S. Mathematics",
			"Portland State University",
			"2016",
			"Portland",
			"Oregon"},
	}

	if !reflect.DeepEqual(fields, expected) {
		t.Errorf("%v: unexpected fields", fields)
	}
}

func TestParseReview(t *testing.T) {
	text, err := ioutil.ReadFile(filepath.Join(testdata, "review.txt"))
	if err != nil {
		t.Error(err)
	}

	fields, err := ParseLines(string(text), []Key{
		{Type: "single", Name: "id"},
		{Type: "single", Name: "drug"},
		{Type: "single", Name: "condition"},
		{Type: "single", Name: "review"},
		{Type: "single", Name: "rating"},
		{Type: "single", Name: "date"},
		{Type: "single", Name: "likes"},
	})
	if err != nil {
		t.Error(err)
	}

	expected := map[string][]string{
		"id":        {"81653"},
		"drug":      {"Liraglutide"},
		"condition": {"Obesity"},
		"rating":    {"8.0"},
		"review":    {"\"I just started the medication this week. I&#039;ve had three shots so far.  I was really worried about the needles but you honestly don&#039;t feel a thing.  I didn&#039;t even feel a pinch when I did the injection.   I do feel a bit tired. I have a very dull headache.     I do have a little bit of  nausea.   I have been injecting in my upper thigh where there&#039;s plenty of fat. (Lol).   My appetite  has decreased so I have to keep an eye on my emotional eating.  I am a former gastric bypass patient  trying to lose the last 30 pounds. The first hundred was easy. So far so good.\""},
		"likes":     {"19"},
		"date":      {"June 24, 2017"},
	}

	if !reflect.DeepEqual(fields, expected) {
		t.Errorf("%v: unexpected fields", fields)
	}
}

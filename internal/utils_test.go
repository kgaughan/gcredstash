package gcredstash_test

import (
	"testing"

	. "github.com/kgaughan/gcredstash/internal"
)

func TestAtoi(t *testing.T) {
	expected := 100
	actual := Atoi("100")

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestVersionNumToStr(t *testing.T) {
	expected := "0000000000000000001"
	actual := VersionNumToStr(1)

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestMapToJSON(t *testing.T) {
	m := map[string]string{"foo": "bar", "bar": "zoo"}

	expected := `{
  "bar": "zoo",
  "foo": "bar"
}`

	actual := MapToJSON(m)

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestMapToJSONWithoutEscape(t *testing.T) {
	m := map[string]string{"<foo>": "&bar", "&bar": "<zoo>"}

	expected := `{
  "&bar": "<zoo>",
  "<foo>": "&bar"
}`

	actual := MapToJSON(m)

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestMaxKeyLen(t *testing.T) {
	key1 := "12"
	val1 := "foobar"
	key2 := "123"
	val2 := "barbaz"

	m := map[*string]*string{&key1: &val1, &key2: &val2}
	expected := 3
	actual := MaxKeyLen(m)

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}
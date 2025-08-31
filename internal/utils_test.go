package internal

import (
	"reflect"
	"testing"
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
}
`

	actual, err := JSONMarshal(m)
	if err != nil {
		t.Error(err)
	}

	if expected != string(actual) {
		t.Errorf("\nexpected: %q\ngot: %q\n", expected, string(actual))
	}
}

func TestMaxKeyLen(t *testing.T) {
	key1 := "12"
	val1 := "foobar"
	key2 := "123"
	val2 := "barbaz"

	m := map[string]string{key1: val1, key2: val2}
	expected := 3
	actual := MaxKeyLen(m)

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestParseContext(t *testing.T) {
	args := []string{"foo=100", "bar=ZOO"}
	expected := map[string]string{"foo": "100", "bar": "ZOO"}
	actual, err := ParseContext(args)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestErrParseContext1(t *testing.T) {
	args := []string{"foo=100", "bar"}
	expected := `invalid context: "bar"`
	if _, err := ParseContext(args); err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

func TestErrParseContext2(t *testing.T) {
	args := []string{"foo=100", "bar="}
	expected := `invalid context: "bar="`
	if _, err := ParseContext(args); err == nil || err.Error() != expected {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, err)
	}
}

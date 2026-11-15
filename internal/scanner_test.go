package scanner

import (
	"bufio"
	"strings"
	"testing"
	"maps"
	"slices"
)

func TestBlockComment(t *testing.T) {
	txt := "/*this is another\n * block comment\n*/"
	want := Token{
		tokType: TokenTypeComment,
		lexeme:  "this is another\n * block comment\n",
		line:    1,
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	got, err := scanner.Scan()

	if err != nil || *got != want {
		t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
	}
}

func TestComment(t *testing.T) {
	txt := "// this is another\n * block comment\n*/"
	want := Token{
		tokType: TokenTypeComment,
		lexeme:  " this is another",
		line:    1,
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	got, err := scanner.Scan()

	if err != nil || *got != want {
		t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
	}
}

func TestMultilineString(t *testing.T) {
	txt := "\"this is another\n string\n\""
	want := Token{
		tokType: TokenTypeString,
		lexeme:  "this is another\n string\n",
		line:    1,
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	got, err := scanner.Scan()

	if err != nil || *got != want {
		t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
	}
}

func TestDecimal(t *testing.T) {
	txt := ".12312\n12"
	expected := []Token{
		Token{
			tokType: TokenTypeNumber,
			lexeme:  ".12312",
			line:    2,
		},
		Token{
			tokType: TokenTypeNumber,
			lexeme:  "12",
			line:    3,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	tokens, err := scanner.ScanAll()

	for i, got := range tokens {
		want := expected[i]
		if err != nil || *got != want {
			t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
		}
	}
}

func TestNumbers(t *testing.T) {
	txt := "124.344\n.12312\n12"
	expected := []Token{
		Token{
			tokType: TokenTypeNumber,
			lexeme:  "124.344",
			line:    1,
		},
		Token{
			tokType: TokenTypeNumber,
			lexeme:  ".12312",
			line:    2,
		},
		Token{
			tokType: TokenTypeNumber,
			lexeme:  "12",
			line:    3,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	tokens, err := scanner.ScanAll()

	for i, got := range tokens {
		want := expected[i]
		if err != nil || *got != want {
			t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
		}
	}
}


func TestIdentifiers(t *testing.T) {
	txt := "variableNzme\n_anotherones"
	expected := []Token{
		Token{
			tokType: TokenTypeIdentifier,
			lexeme:  "variableNzme",
			line:    1,
		},
		Token{
			tokType: TokenTypeIdentifier,
			lexeme:  "_anotherone",
			line:    2,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	tokens, err := scanner.ScanAll()

	for i, got := range tokens {
		want := expected[i]
		if err != nil || *got != want {
			t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
		}
	}
}


func TestKeywords(t *testing.T) {
	kws := slices.Collect(maps.Keys(keywords))
	txt := strings.Join(kws, "\n")
	r := bufio.NewReader(strings.NewReader(txt))
	scanner := NewScanner(r)

	tokens, err := scanner.ScanAll()

	for i, got := range tokens {
		want := Token{
			tokType: keywords[kws[i]],
			lexeme:  kws[i],
			line:    i,
		}

		if err != nil || *got != want {
			t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
		}
	}
}

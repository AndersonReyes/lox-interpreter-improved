package scanner

import (
	"bufio"
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
)

func removeIgnoredTokens(tokens []*Token) []*Token {
	out := []*Token{}

	for _, t := range tokens {
		if t.tokType == TokenTypeIgnore {
			continue
		}
		out = append(out, t)
	}

	return out
}

func expectTokens(t *testing.T, got []*Token, want []Token) {
	got = removeIgnoredTokens(got)

	if len(got) != len(want) {
		for _, tok := range got {
			fmt.Printf("got: %+v\n", *tok)
		}
		t.Errorf("got %d more tokens than expected", len(got))
	}

	for i, g := range got {
		w := want[i]
		if *g != w {
			t.Errorf("\ngot:\n%+v, \n\nbut wanted:\n%+v", g, w)
		}
	}
}

func TestBlockComment(t *testing.T) {
	txt := "/*this is another\n * block comment\n*/"
	want := Token{
		tokType: TokenTypeComment,
		lexeme:  "this is another\n * block comment\n",
		line:    1,
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	got, err := scanner.Scan()

	if err != nil || *got != want {
		t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
	}
}

func TestComment(t *testing.T) {
	txt := "// this is a comment"
	want := Token{
		tokType: TokenTypeComment,
		lexeme:  " this is a comment",
		line:    1,
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

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
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	got, err := scanner.Scan()

	if err != nil || *got != want {
		t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
	}
}

func TestDecimal(t *testing.T) {
	txt := ".12312\n12"
	expected := []Token{
		{
			tokType: TokenTypeNumber,
			lexeme:  ".12312",
			line:    1,
		},
		{
			tokType: TokenTypeNumber,
			lexeme:  "12",
			line:    2,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	tokens, err := scanner.ScanAll()
	if err != nil {
		t.Errorf("error scanning tokens: %+v", err)
	}

	expectTokens(t, tokens, expected)

}

func TestNumbers(t *testing.T) {
	txt := "124.344\n.12312\n12"
	expected := []Token{
		{
			tokType: TokenTypeNumber,
			lexeme:  "124.344",
			line:    1,
		},
		{
			tokType: TokenTypeNumber,
			lexeme:  ".12312",
			line:    2,
		},
		{
			tokType: TokenTypeNumber,
			lexeme:  "12",
			line:    3,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	tokens, err := scanner.ScanAll()
	if err != nil {
		t.Errorf("error scanning tokens: %+v", err)
	}

	expectTokens(t, tokens, expected)
}

func TestIdentifiers(t *testing.T) {
	txt := "variableNzme\n_anotherones"
	expected := []Token{
		{
			tokType: TokenTypeIdentifier,
			lexeme:  "variableNzme",
			line:    1,
		},
		{
			tokType: TokenTypeIdentifier,
			lexeme:  "_anotherones",
			line:    2,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	tokens, err := scanner.ScanAll()
	if err != nil {
		t.Errorf("error scanning tokens: %+v", err)
	}

	expectTokens(t, tokens, expected)
}

func TestKeywords(t *testing.T) {
	kws := slices.Collect(maps.Keys(keywords))
	txt := strings.Join(kws, "\n")
	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	tokens, err := scanner.ScanAll()
	if err != nil {
		t.Errorf("error scanning tokens: %+v", err)
	}

	tokens = removeIgnoredTokens(tokens)

	for i, got := range tokens {
		want := Token{
			tokType: keywords[kws[i]],
			lexeme:  kws[i],
			line:    i + 1,
		}

		if err != nil || *got != want {
			t.Errorf("\ngot:\n%+v, %+v,\n\nbut wanted:\n%+v", got, err, want)
		}
	}
}

func TestManyStatements(t *testing.T) {
	txt := strings.Join([]string{
		"\"this is\na string\"",
		"124.344",
		".12312",
		"// this is a comment",
		"/* This is a nother\n* block comment\n*/",
	}, "\n")

	expected := []Token{
		{
			tokType: TokenTypeString,
			lexeme:  "this is\na string",
			line:    1,
		},
		{
			tokType: TokenTypeNumber,
			lexeme:  "124.344",
			line:    2,
		},
		{
			tokType: TokenTypeNumber,
			lexeme:  ".12312",
			line:    3,
		},
		{
			tokType: TokenTypeComment,
			lexeme:  " this is a comment",
			line:    4,
		},
		{
			tokType: TokenTypeComment,
			lexeme:  " This is a nother\n* block comment\n",
			line:    5,
		},
	}

	r := bufio.NewReader(strings.NewReader(txt))
	scanner, err := NewScanner(r)
	if err != nil {
		t.Errorf("failed to init scanner: %s", err)

	}

	tokens, err := scanner.ScanAll()

	if err != nil {
		t.Errorf("error scanning tokens: %+v", err)
	}

	expectTokens(t, tokens, expected)

}

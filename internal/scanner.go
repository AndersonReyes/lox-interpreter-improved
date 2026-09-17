package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
)

type TokenType int

const (
	TokenTypeLeftParen = iota
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeSlash
	TokenTypeStar

	TokenTypeBang         // !
	TokenTypeBangEqual    // !=
	TokenTypeEqual        // =
	TokenTypeEqualEqual   // ==
	TokenTypeGreater      // >
	TokenTypeGreaterEqual // >=
	TokenTypeLess         // <
	TokenTypeLessEqual    // <=

	TokenTypeIdentifier   // var name
	TokenTypeString       // string type
	TokenTypeNumber       // number type
	TokenTypeQuestionMark // ? operator for option types

	// KEYWORDS
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFun
	TokenTypeFor
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile
	TokenTypeEof

	// to detect errors / invalid tokens
	TokenTypeUnrecognized
)

type Token struct {
	tokType TokenType
	lexeme  string
	line    int
}

type Scanner struct {
	reader *bufio.Reader
}

func NewScanner(reader *bufio.Reader) *Scanner {
	return &Scanner{reader: reader}
}

func (s *Scanner) parseToken(line int, currChar string, nextChar string) (*Token, error) {

	switch currChar {
	case "(":
		return &Token{
			tokType: TokenTypeLeftParen,
			lexeme:  currChar,
			line:    line,
		}, nil
	case ")":
		return &Token{
			tokType: TokenTypeRightParen,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "{":
		return &Token{
			tokType: TokenTypeLeftBrace,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "}":
		return &Token{
			tokType: TokenTypeRightBrace,
			lexeme:  currChar,
			line:    line,
		}, nil
	case ",":
		return &Token{
			tokType: TokenTypeComma,
			lexeme:  currChar,
			line:    line,
		}, nil
	case ".":
		return &Token{
			tokType: TokenTypeDot,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "-":
		return &Token{
			tokType: TokenTypeMinus,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "+":
		return &Token{
			tokType: TokenTypePlus,
			lexeme:  currChar,
			line:    line,
		}, nil
	case ";":
		return &Token{
			tokType: TokenTypeSemicolon,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "*":
		return &Token{
			tokType: TokenTypeStar,
			lexeme:  currChar,
			line:    line,
		}, nil
	case "!":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeBangEqual,
				lexeme:  "!=",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeBang,
				lexeme:  currChar,
				line:    line,
			}, nil
		}
	case "=":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeEqualEqual,
				lexeme:  "==",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeEqual,
				lexeme:  currChar,
				line:    line,
			}, nil
		}
	case "<":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeLessEqual,
				lexeme:  "<=",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeLess,
				lexeme:  currChar,
				line:    line,
			}, nil
		}
	case ">":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeGreaterEqual,
				lexeme:  ">=",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeGreater,
				lexeme:  currChar,
				line:    line,
			}, nil
		}
	default:
		return nil, reportError(line, currChar, "invalid character")
	}
}

func (s *Scanner) consume() (string, string, error) {
	r, _, err := s.reader.ReadRune()

	currChar := string(r)
	if err != nil {
		if err == io.EOF {
			return "", "", err
		} else {
			return "", "", errors.Join(err, fmt.Errorf("invalid token %s", currChar))

		}
	}

	nextRune, _, err := s.reader.ReadRune()
	nextChar := string(nextRune)
	if err != nil {
		if err == io.EOF {
			return currChar, "", nil
		} else {
			return "", "", errors.Join(err, fmt.Errorf("invalid next token %s", nextChar))
		}
	}

	// Move the cursor back since we just want to peek at the next rune
	if err := s.reader.UnreadRune(); err != nil {
		return "", "", errors.Join(err, errors.New("failed to unread rune"))
	}

	return currChar, nextChar, nil
}

func (s *Scanner) Scan() error {
	line := 1
	for {
		currChar, nextChar, err := s.consume()

		if err == io.EOF {
			break
		}

		if currChar == "\n" {
			line += 1
			continue
		}

		// fmt.Printf("Current, next and line %s -> %s at line %d", currChar, nextChar, line)
		tok, err := s.parseToken(line, currChar, nextChar)
		if err != nil {
			return err
		}

		fmt.Printf("Token: %+v\n", tok)

	}

	return nil

}

func reportError(line int, where string, message string) error {
	msg := fmt.Sprintf("[line %d] Error %s: %s", line, where, message)
	log.Print(msg)
	return errors.New(msg)
}

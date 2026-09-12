package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
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
	literal string // TODO can we store the actual type here?
	line    int
}

func getTokenType(ch string) TokenType {
	switch ch {
	case "(":
		return TokenTypeLeftParen
	case ")":
		return TokenTypeRightParen
	case "{":
		return TokenTypeLeftBrace
	case "}":
		return TokenTypeRightBrace
	case ",":
		return TokenTypeComma
	case ".":
		return TokenTypeDot
	case "-":
		return TokenTypeMinus
	case "+":
		return TokenTypePlus
	case ";":
		return TokenTypeSemicolon
	case "*":
		return TokenTypeStar
	default:
		return TokenTypeUnrecognized
	}
}

func ScanTokens(reader *bufio.Reader) error {
	line := 0
	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		ch := string(r)
		if err != nil {
			return reportError(line, "scanner", "invalid token "+ch)
		}

		tokType := getTokenType(ch)
		if tokType == TokenTypeUnrecognized {
			return errors.Join(err, reportError(line, ch, "unrecognized character"))
		}

		token := Token{
			tokType: tokType,
			lexeme:  ch,
			literal: ch,
			line:    line,
		}

		fmt.Printf("Token: %+v\n", token)

		if strings.Compare(ch, "\n") == 0 {
			line += 1
			continue
		}
	}

	return nil

}

func reportError(line int, where string, message string) error {
	msg := fmt.Sprintf("[line %d] Error %s: %s", line, where, message)
	log.Print(msg)
	return errors.New(msg)
}

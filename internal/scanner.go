package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const numberParts = ".0123456789"

var alphanumRe = regexp.MustCompile(`^[_a-zA-Z0-9]`)

type TokenType string

const (
	TokenTypeLeftParen  = "("
	TokenTypeRightParen = ")"
	TokenTypeLeftBrace  = "{"
	TokenTypeRightBrace = "}"
	TokenTypeComma      = "Comma"
	TokenTypeDot        = "Period"
	TokenTypeMinus      = "Minus"
	TokenTypePlus       = "Plus"
	TokenTypeSemicolon  = ";"
	TokenTypeSlash      = "/"
	TokenTypeStar       = "*"

	TokenTypeBang         = "!"
	TokenTypeBangEqual    = "!="
	TokenTypeEqual        = "="
	TokenTypeEqualEqual   = "=="
	TokenTypeGreater      = ">"
	TokenTypeGreaterEqual = ">="
	TokenTypeLess         = "<"
	TokenTypeLessEqual    = "<="

	TokenTypeIdentifier   = "Identifier"
	TokenTypeString       = "TypeString"
	TokenTypeNumber       = "TypeNumber"
	TokenTypeQuestionMark = "?"

	// KEYWORDS
	TokenTypeAnd    = "and"
	TokenTypeClass  = "class"
	TokenTypeElse   = "else"
	TokenTypeFalse  = "false"
	TokenTypeFun    = "fun"
	TokenTypeFor    = "for"
	TokenTypeIf     = "if"
	TokenTypeNil    = "nil"
	TokenTypeOr     = "or"
	TokenTypePrint  = "print"
	TokenTypeReturn = "return"
	TokenTypeSuper  = "super"
	TokenTypeThis   = "this"
	TokenTypeTrue   = "true"
	TokenTypeVar    = "var"
	TokenTypeWhile  = "while"

	// to detect errors / invalid tokens
	TokenTypeUnrecognized = "Unknown"

	TokenTypeComment = "Comment"
	TokenTypeIgnore  = "IgnoreThisChar"
	TokenTypeEof     = "Eof"
)

var keywords = map[string]TokenType{
	TokenTypeAnd:    TokenTypeAnd,
	TokenTypeClass:  TokenTypeClass,
	TokenTypeElse:   TokenTypeElse,
	TokenTypeFalse:  TokenTypeFalse,
	TokenTypeFun:    TokenTypeFun,
	TokenTypeFor:    TokenTypeFor,
	TokenTypeIf:     TokenTypeIf,
	TokenTypeNil:    TokenTypeNil,
	TokenTypeOr:     TokenTypeOr,
	TokenTypePrint:  TokenTypePrint,
	TokenTypeReturn: TokenTypeReturn,
	TokenTypeSuper:  TokenTypeSuper,
	TokenTypeThis:   TokenTypeThis,
	TokenTypeTrue:   TokenTypeTrue,
	TokenTypeVar:    TokenTypeVar,
	TokenTypeWhile:  TokenTypeWhile,
}

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
	case " ", "\r", "\t":
		return &Token{
			tokType: TokenTypeIgnore,
			lexeme:  currChar,
			line:    line,
		}, nil

	case "(":
		return &Token{
			tokType: TokenTypeLeftParen,
			lexeme:  "",
			line:    line,
		}, nil
	case ")":
		return &Token{
			tokType: TokenTypeRightParen,
			lexeme:  "",
			line:    line,
		}, nil
	case "{":
		return &Token{
			tokType: TokenTypeLeftBrace,
			lexeme:  "",
			line:    line,
		}, nil
	case "}":
		return &Token{
			tokType: TokenTypeRightBrace,
			lexeme:  "",
			line:    line,
		}, nil
	case ",":
		return &Token{
			tokType: TokenTypeComma,
			lexeme:  ",",
			line:    line,
		}, nil
	case ".":
		return &Token{
			tokType: TokenTypeDot,
			lexeme:  ".",
			line:    line,
		}, nil
	case "-":
		return &Token{
			tokType: TokenTypeMinus,
			lexeme:  "-",
			line:    line,
		}, nil
	case "+":
		return &Token{
			tokType: TokenTypePlus,
			lexeme:  "+",
			line:    line,
		}, nil
	case ";":
		return &Token{
			tokType: TokenTypeSemicolon,
			lexeme:  "",
			line:    line,
		}, nil
	case "*":
		return &Token{
			tokType: TokenTypeStar,
			lexeme:  "",
			line:    line,
		}, nil
	case "!":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeBangEqual,
				lexeme:  "",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeBang,
				lexeme:  "",
				line:    line,
			}, nil
		}
	case "=":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeEqualEqual,
				lexeme:  "",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeEqual,
				lexeme:  "",
				line:    line,
			}, nil
		}
	case "<":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeLessEqual,
				lexeme:  "",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeLess,
				lexeme:  "",
				line:    line,
			}, nil
		}
	case ">":
		if nextChar == "=" {
			// consume nextChar
			_, _, _ = s.consume()
			return &Token{
				tokType: TokenTypeGreaterEqual,
				lexeme:  "",
				line:    line,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeGreater,
				lexeme:  "",
				line:    line,
			}, nil
		}
	case "/":
		// Handle comments
		if nextChar == "/" {
			// this is a comment read until end of line
			builder := strings.Builder{}

			for {
				c, _, err := s.consume()
				if c == "\n" || err == io.EOF {
					break
				}

				if err != nil {
					return nil, reportError(line, currChar, "invalid comment character")
				}

				builder.WriteString(c)
			}

			return &Token{
				tokType: TokenTypeComment,
				lexeme:  builder.String(),
				line:    line,
			}, nil

		} else {
			return &Token{
				tokType: TokenTypeSlash,
				lexeme:  "",
				line:    line,
			}, nil
		}
	case "\"": // String literals
		stringLit, err := s.parseString(line)
		if err != nil {
			return nil, err
		}

		return &Token{
			tokType: TokenTypeString,
			lexeme:  stringLit,
			line:    line,
		}, nil
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		num, err := s.parseNumber(line)
		if err != nil {
			return nil, err
		}

		return &Token{
			tokType: TokenTypeNumber,
			lexeme:  num,
			line:    line,
		}, nil
	default:
		if alphanumRe.MatchString(currChar) {
			builder := strings.Builder{}
			builder.WriteString(currChar)

			for {
				c, _, err := s.consume()

				if err != nil {
					return nil, reportError(line, currChar, "error parsing Identifier")
				}

				// log.Printf("default: builder=%s, %s=%t", builder.String(), c, alphanumRe.MatchString(c))

				if !alphanumRe.MatchString(c) {
					break
				}

				builder.WriteString(c)
			}
			lexeme := builder.String()

			tokType, is_keyword := keywords[lexeme]
			if !is_keyword {
				tokType = TokenTypeIdentifier
			}

			return &Token{
				tokType: tokType,
				lexeme:  lexeme,
				line:    line,
			}, nil
		} else {
			return nil, reportError(line, currChar, "invalid character")
		}
	}
}

func (s *Scanner) parseNumber(line int) (string, error) {
	builder := strings.Builder{}
	alreadyHasPeriod := false
	for {
		c, _, err := s.consume()

		if !strings.Contains(numberParts, c) || err == io.EOF {
			break
		}

		if c == "." {
			if alreadyHasPeriod {
				return "", reportError(line, builder.String(), "invalid number, duplicate periods")
			} else {
				alreadyHasPeriod = true
			}
		}

		builder.WriteString(c)
	}

	numStr := builder.String()
	// TODO: figure out how to return numbers or strings inside the struct
	_, err := strconv.ParseFloat(numStr, 64)

	if err != nil {
		return "", reportError(line, numStr, "invalid number")
	}

	if strings.HasSuffix(numStr, ".") {
		return "", reportError(line, builder.String(), "number ends with a period")
	}
	return numStr, nil
}

func (s *Scanner) parseString(line int) (string, error) {
	builder := strings.Builder{}
	for {
		c, _, err := s.consume()
		if err == io.EOF {
			return "", reportError(line, builder.String(), "unterminated string.")
		}

		if c == "\"" {
			break
		}
		builder.WriteString(c)
	}

	return builder.String(), nil
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
	// TODO: store the next char or string in *Scanner. This avoids callings unreadrune a few lines down
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
	msg := fmt.Sprintf("[line %d] Error near \"%s\": %s", line, where, message)
	log.Print(msg)
	return errors.New(msg)
}

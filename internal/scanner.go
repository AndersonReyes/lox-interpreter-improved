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
	TokenTypeNewLine = "\\n"
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
	reader        *bufio.Reader
	currentLine   int
	currentColumn int
	currChar      string
	nextChar      string
}

func NewScanner(reader *bufio.Reader) *Scanner {
	log.Printf("")
	s := &Scanner{reader: reader, currentLine: 1}
	// move the scanner to populate curr and next chars for the first time
	_ = s.consume()
	_ = s.consume()
	return s
}

func (s *Scanner) parseToken() (*Token, error) {

	switch s.currChar {
	case "", " ", "\r", "\t", "\n":
		tok := &Token{
			tokType: TokenTypeIgnore,
			lexeme:  s.currChar,
			line:    s.currentLine,
		}

		// consume the unwanted chars
		err := s.consume()
		if s.currChar == "\n" {
			s.currentLine += 1
			s.currentColumn = 1
		}

		return tok, err

	case "(":
		return &Token{
			tokType: TokenTypeLeftParen,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case ")":
		return &Token{
			tokType: TokenTypeRightParen,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "{":
		return &Token{
			tokType: TokenTypeLeftBrace,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "}":
		return &Token{
			tokType: TokenTypeRightBrace,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case ",":
		return &Token{
			tokType: TokenTypeComma,
			lexeme:  ",",
			line:    s.currentLine,
		}, nil
	case ".":
		if strings.Contains(numberParts, s.nextChar) {
			num, err := s.parseNumber(s.currChar)
			if err != nil {
				return nil, err
			}

			return &Token{
				tokType: TokenTypeNumber,
				lexeme:  num,
				line:    s.currentLine,
			}, nil

		} else {
			return &Token{
				tokType: TokenTypeDot,
				lexeme:  ".",
				line:    s.currentLine,
			}, nil

		}
	case "-":
		return &Token{
			tokType: TokenTypeMinus,
			lexeme:  "-",
			line:    s.currentLine,
		}, nil
	case "+":
		return &Token{
			tokType: TokenTypePlus,
			lexeme:  "+",
			line:    s.currentLine,
		}, nil
	case ";":
		return &Token{
			tokType: TokenTypeSemicolon,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "*":
		return &Token{
			tokType: TokenTypeStar,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "!":
		if s.nextChar == "=" {
			// consume nextChar
			_ = s.consume()
			return &Token{
				tokType: TokenTypeBangEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeBang,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "=":
		if s.nextChar == "=" {
			// consume nextChar
			_ = s.consume()
			return &Token{
				tokType: TokenTypeEqualEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "<":
		if s.nextChar == "=" {
			// consume nextChar
			_ = s.consume()
			return &Token{
				tokType: TokenTypeLessEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeLess,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case ">":
		if s.nextChar == "=" {
			// consume nextChar
			_ = s.consume()
			return &Token{
				tokType: TokenTypeGreaterEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			return &Token{
				tokType: TokenTypeGreater,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "/":
		// Handle comments
		if s.nextChar == "/" || s.nextChar == "*" {
			// skip the slash or star for comment
			isBlockComment := false
			if s.nextChar == "*" {
				isBlockComment = true
			}

			_ = s.consume()

			// this is a comment read until end of line
			builder := strings.Builder{}

			for {
				err := s.consume()
				// log.Printf("comment] c={%s}, nextcurrent={%s}, line: %d, err=%+v, is block comment=%v\n", s.currChar, s.nextChar, s.currentLine, err, isBlockComment)

				// check if we are at the end, or if it's a block comment end or new line if its not a block comment
				if !isBlockComment && s.currChar == "\n" {
					break
				}

				if isBlockComment && s.currChar == "*" && s.nextChar == "/" {
					// skep the next two chars since we are already including it in the builder above
					_ = s.consume()
					break
				}

				if err != nil {
					return nil, reportError(s.currentLine,s.currentColumn, s.currChar, "invalid comment character")
				}

				builder.WriteString(s.currChar)
			}

			return &Token{
				tokType: TokenTypeComment,
				lexeme:  builder.String(),
				line:    s.currentLine,
			}, nil

		} else {
			return &Token{
				tokType: TokenTypeSlash,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "\"": // String literals
		stringLit, err := s.parseString()
		if err != nil {
			return nil, err
		}

		return &Token{
			tokType: TokenTypeString,
			lexeme:  stringLit,
			line:    s.currentLine,
		}, nil
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		num, err := s.parseNumber(s.currChar)
		if err != nil {
			return nil, err
		}

		return &Token{
			tokType: TokenTypeNumber,
			lexeme:  num,
			line:    s.currentLine,
		}, nil
	default:
		if alphanumRe.MatchString(s.currChar) {
			builder := strings.Builder{}
			builder.WriteString(s.currChar)

			for {
				err := s.consume()

				if err != nil && err != io.EOF {
					return nil, reportError(s.currentLine, s.currentColumn, s.currChar, "error parsing Identifier")
				}

				// log.Printf("default: builder=%s, %s=%t", builder.String(), c, alphanumRe.MatchString(c))

				if !alphanumRe.MatchString(s.currChar) {
					break
				}

				builder.WriteString(s.currChar)

				if err == io.EOF {
					break
				}
			}
			lexeme := builder.String()

			tokType, is_keyword := keywords[lexeme]
			if !is_keyword {
				tokType = TokenTypeIdentifier
			}

			return &Token{
				tokType: tokType,
				lexeme:  lexeme,
				line:    s.currentLine,
			}, nil
		} else {
			return nil, reportError(s.currentLine, s.currentColumn, s.currChar, "invalid character")
		}
	}
}

func (s *Scanner) parseNumber(prefix string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(prefix)
	alreadyHasPeriod := strings.Contains(prefix, ".")
	for {
		err := s.consume()

		if !strings.Contains(numberParts, s.currChar) {
			break
		}

		if err != nil && err != io.EOF {
			return "", err
		}

		if s.currChar == "." {
			if alreadyHasPeriod {
				return "", reportError(s.currentLine, s.currentColumn, builder.String(), "invalid number, duplicate periods")
			} else {
				alreadyHasPeriod = true
			}
		}

		builder.WriteString(s.currChar)

		if err == io.EOF {
			break
		}
	}

	numStr := builder.String()
	// TODO: figure out how to return numbers or strings inside the struct
	_, err := strconv.ParseFloat(numStr, 64)

	if err != nil {
		return "", reportError(s.currentLine, s.currentColumn, numStr, "invalid number")
	}

	if strings.HasSuffix(numStr, ".") {
		return "", reportError(s.currentLine, s.currentColumn, builder.String(), "number ends with a period")
	}
	return numStr, nil
}

func (s *Scanner) parseString() (string, error) {
	builder := strings.Builder{}
	for {
		err := s.consume()
		if err != nil && err != io.EOF {
			return "",errors.Join(err,  reportError(s.currentLine, s.currentColumn, builder.String(), "invalid string."))
		}
		//log.Printf("string] c={%s}, nextcurrent={%s}, line: %d\n", s.currChar, s.nextChar, s.currentLine)

		if s.currChar == "\"" {
			// consume closing quote
			_ = s.consume()
			break
		}

		builder.WriteString(s.currChar)

		if err == io.EOF {
			return "",errors.Join(err,  reportError(s.currentLine, s.currentColumn, builder.String(), "unterminated string."))
		}
	}

	return builder.String(), nil
}

func (s *Scanner) consume() error {
	r, _, err := s.reader.ReadRune()

	s.currChar = s.nextChar
	s.nextChar = string(r)
	s.currentColumn += 1
	if err != nil {
		if err == io.EOF {
			return err
		} else {
			return errors.Join(err, fmt.Errorf("invalid token %s, next=%s", s.currChar, s.nextChar))

		}
	}

	return nil
}

func (s *Scanner) Scan() (*Token, error) {
	tok, err := s.parseToken()
	if err != nil {
		return nil, err
	}

	return tok, nil

}

func (s *Scanner) ScanAll() ([]*Token, error) {
	tokens := []*Token{}
	for {
		tok, err := s.Scan()

		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			return nil, err
		}

		tokens = append(tokens, tok)
	}
	return tokens, nil

}

func reportError(line int, column int, where string, message string) error {
	msg := fmt.Sprintf("[line %d, column %d] Error near \"%s\": %s", line, column, where, message)
	//log.Print(msg)
	return errors.New(msg)
}

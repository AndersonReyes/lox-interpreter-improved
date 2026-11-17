package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

var alphanumRe = regexp.MustCompile(`^[_a-zA-Z0-9]`)
var digitOrPeriod = regexp.MustCompile(`^[0-9\.]`)

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
	currentLine   int
	currentColumn int
	state         *ScannerState
}

func NewScanner(reader *bufio.Reader) (*Scanner, error) {
	_ = log.Flags()
	state, err := NewScannerState(reader)

	if err != nil {
		return nil, err
	}
	s := &Scanner{state: state, currentLine: 1}
	return s, nil
}

func (s *Scanner) parseToken() (*Token, error) {

	switch s.state.Head() {
	case "", " ", "\r", "\t", "\n":
		tok := &Token{
			tokType: TokenTypeIgnore,
			lexeme:  s.state.Head(),
			line:    s.currentLine,
		}

		if s.state.Head() == "\n" {
			s.currentLine += 1
			s.currentColumn = 1
		}

		s.state.Skip(1)

		return tok, nil

	case "(":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeLeftParen,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case ")":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeRightParen,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "{":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeLeftBrace,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "}":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeRightBrace,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case ",":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeComma,
			lexeme:  ",",
			line:    s.currentLine,
		}, nil
	case ".":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeDot,
			lexeme:  ".",
			line:    s.currentLine,
		}, nil
	case "-":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeMinus,
			lexeme:  "-",
			line:    s.currentLine,
		}, nil
	case "+":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypePlus,
			lexeme:  "+",
			line:    s.currentLine,
		}, nil
	case ";":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeSemicolon,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "*":
		s.state.Skip(1)
		return &Token{
			tokType: TokenTypeStar,
			lexeme:  "",
			line:    s.currentLine,
		}, nil
	case "!":
		if s.state.Peek() == "=" {
			s.state.Skip(2)
			return &Token{
				tokType: TokenTypeBangEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			s.state.Skip(1)
			return &Token{
				tokType: TokenTypeBang,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "=":
		if s.state.Peek() == "=" {
			s.state.Skip(2)
			return &Token{
				tokType: TokenTypeEqualEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			s.state.Skip(1)
			return &Token{
				tokType: TokenTypeEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "<":
		if s.state.Peek() == "=" {
			s.state.Skip(2)
			return &Token{
				tokType: TokenTypeLessEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			s.state.Skip(1)
			return &Token{
				tokType: TokenTypeLess,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case ">":
		if s.state.Peek() == "=" {
			s.state.Skip(2)
			return &Token{
				tokType: TokenTypeGreaterEqual,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		} else {
			s.state.Skip(1)
			return &Token{
				tokType: TokenTypeGreater,
				lexeme:  "",
				line:    s.currentLine,
			}, nil
		}
	case "/":
		// Handle comments
		if s.state.Peek() == "/" || s.state.Peek() == "*" {
			isBlock := s.state.Peek() == "*"
			s.state.Skip(2)
			block, err := s.parseComment(isBlock)
			if err != nil {
				return nil, err
			}
			return &Token{
				tokType: TokenTypeComment,
				lexeme:  block,
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
		num, err := s.parseNumber()
		if err != nil {
			return nil, err
		}

		return &Token{
			tokType: TokenTypeNumber,
			lexeme:  num,
			line:    s.currentLine,
		}, nil
	default:
		lexeme, err := s.parseIdentifier()
		if err != nil {
			return nil, err
		}

		tokType, is_keyword := keywords[lexeme]
		if !is_keyword {
			tokType = TokenTypeIdentifier
		}

		return &Token{
			tokType: tokType,
			lexeme:  lexeme,
			line:    s.currentLine,
		}, nil
	}
}

func (s *Scanner) parseComment(isBlock bool) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(s.state.Head())
	for {
		c, err := s.state.GetNext()

		if err != nil && err != io.EOF {
			return "", errors.Join(err, reportError(s.currentLine, s.currentColumn, builder.String(), "error parsing comment"))
		}

		if isBlock && c == "*" && s.state.Peek() == "/" {
			s.state.Skip(2)
			break
		}

		if !isBlock && c == "\n" {
			break
		}

		builder.WriteString(c)

		if err == io.EOF {
			break
		}

	}

	return builder.String(), nil
}

func (s *Scanner) parseIdentifier() (string, error) {
	builder := strings.Builder{}
	builder.WriteString(s.state.Head())

	for {
		c, err := s.state.GetNext()
		if err != nil {
			return "", reportError(s.currentLine, s.currentColumn, builder.String(), "error parsing identifier")
		}

		if !alphanumRe.MatchString(c) {
			break
		}

		builder.WriteString(c)
	}

	matched := builder.String()

	if digitOrPeriod.MatchString(matched) {
		return "", reportError(s.currentLine, s.currentColumn, matched, "invalid start of identifier")
	}

	return matched, nil

}

func (s *Scanner) parseNumber() (string, error) {
	builder := strings.Builder{}
	builder.WriteString(s.state.Head())

	alreadyHasPeriod := s.state.Head() == "."

	for {
		c, err := s.state.GetNext()

		if err != nil {
			return "", errors.Join(err, reportError(s.currentLine, s.currentColumn, builder.String(), "error parsing number"))
		}

		if !digitOrPeriod.MatchString(c) {
			break
		}

		if c == "." {
			if alreadyHasPeriod {
				return "", reportError(s.currentLine, s.currentColumn, builder.String(), "invalid number, duplicate periods")
			} else {
				alreadyHasPeriod = true
			}
		}

		builder.WriteString(c)
	}

	// log.Printf("parse number: c %s, n %s, numstr %s", s.currChar, s.nextChar, numStr)
	return builder.String(), nil
}

func (s *Scanner) parseString() (string, error) {
	builder := strings.Builder{}

	for {
		c, err := s.state.GetNext()

		if err == io.EOF {
			return "", errors.Join(err, reportError(s.currentLine, s.currentColumn, builder.String(), "unterminated string."))
		}

		if err != nil {
			return "", errors.Join(err, reportError(s.currentLine, s.currentColumn, builder.String(), "error parsing string"))
		}

		//log.Printf("string] c={%s}, nextcurrent={%s}, line: %d\n", s.currChar, s.nextChar, s.currentLine)

		if c == "\"" {
			break
		}

		builder.WriteString(c)
	}

	return builder.String(), nil
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

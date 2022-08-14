package bnf

import (
	"fmt"
	"io"
	"strings"
)

// A Token is an element of BNF grammar; it's any single thing we might
// parse
type Token struct {
	Type  int
	Value string
}

// A TokenStream is a set of tokens
type TokenStream struct {
	buf []Token
}

const (
	TokenComment = iota
	TokenTerminal
	TokenNonterminal
	TokenOpEqual
	TokenOpBar
	TokenEOL
)

func Tokenize(r io.Reader) (*TokenStream, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return tokenizeString(string(buf))
}

func tokenizeString(s string) (*TokenStream, error) {
	var (
		stream = new(TokenStream)
		pos    = 0
		end    = len(s)
	)

	// Nothing to do
	if end == 0 {
		return stream, nil
	}

	for pos < end {
		switch s[pos] {
		case ':':
			if pos+3 >= len(s) {
				return stream, fmt.Errorf("expected '::=', but reached end of input")
			}

			if s[pos:pos+3] != "::=" {
				return stream, fmt.Errorf("expected '::=', got '%s...'", s[pos:3])
			}

			stream.Push(Token{Type: TokenOpEqual})
			pos += 3

		case '|':
			stream.Push(Token{Type: TokenOpBar})
			pos++

		case '"', '\'':
			// A terminal symbol is enclosed in some form of quotation
			// mark
			val, nextPos, err := until(s, pos+1, s[pos])
			if err != nil {
				return stream, err
			}

			stream.Push(Token{Type: TokenTerminal, Value: val})
			pos = nextPos

		case '<':
			// Define a nonterminal symbol as '<...>'
			val, nextPos, err := until(s, pos+1, '>')
			if err != nil {
				return stream, err
			}

			stream.Push(Token{Type: TokenNonterminal, Value: val})
			pos = nextPos

		case ' ', '\t', '\r':
			// Whitespace should be skipped
			pos++

		case '\n':
			// Newlines are significant, and have their own token
			stream.Push(Token{Type: TokenEOL})
			pos++

		default:
			return stream, fmt.Errorf("unexpected character '%c'", s[pos])
		}
	}

	return stream, nil
}

// until returns some encapsulated string (everything _until_ some
// delimiter). It has support for escaped characters to allow you to
// contain the delimiter within the encapsulated form with an escape.
func until(s string, pos int, delim byte) (string, int, error) {
	var (
		buf strings.Builder
		cur int
	)

	for cur = pos; cur < len(s); cur++ {
		// What if we found an escaped character?
		if s[cur] == '\\' && cur+1 < len(s) {
			// We should write everything up until the backslash, and
			// then write the escaped character
			buf.WriteString(s[pos:cur] + string(s[cur+1]))

			// This is annoyingly tricky. We want to set cur to the
			// character that's being escaped. We then want to set pos
			// to the one after the escaped character. And then, when we
			// continue, cur will get to the one after the escaped
			// character as well.
			cur = cur + 1
			pos = cur + 1
			continue
		}

		if s[cur] == delim {
			break
		}
	}

	// Whatever's left, write that into buf
	buf.WriteString(s[pos:cur])

	return buf.String(), cur + 1, nil
}

// Push will add a token to the token buffer
func (s *TokenStream) Push(t Token) {
	s.buf = append(s.buf, t)
}
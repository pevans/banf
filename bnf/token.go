package bnf

import (
	"io"
	"strings"

	"github.com/pevans/banf/position"
	"github.com/rotisserie/eris"
)

// A Token is an element of BNF grammar; it's any single thing we might
// parse
type Token struct {
	Type  int
	Value string
}

// A TokenStream is a set of tokens
type TokenStream struct {
	buf []*Token
	pos int
}

const (
	TokenT = iota
	TokenNT
	TokenEq
	TokenBar
	TokenEOL
)

// Tokenize will take any input (accessible from a Reader) and produce a
// token stream. If an error occurs while parsing, it will return a
// partial token stream.
func Tokenize(r io.Reader) (*TokenStream, *ParseError) {
	stream := new(TokenStream)

	buf, err := io.ReadAll(r)
	if err != nil {
		return stream, &ParseError{Err: err}
	}

	perr := tokenizeString(string(buf), stream)
	return stream, perr
}

// tokenizeString will produce a stream of tokens from a given string.
// If an error is occurred, it will produce a partial token stream that
// contains everything until the error.
func tokenizeString(s string, stream *TokenStream) *ParseError {
	var (
		pos  = 0
		end  = len(s)
		line int
	)

	// Nothing to do
	if end == 0 {
		return nil
	}

	for pos < end {
		switch s[pos] {
		case '#':
			// Skip all characters until the end of the line
			for pos < end && s[pos] != '\n' {
				pos++
			}

		case ':':
			if pos+3 >= len(s) {
				return &ParseError{
					Line:      line,
					Incidence: position.Show(s, pos),
					Err:       eris.New("expected '::=', but reached end of input"),
				}
			}

			if s[pos:pos+3] != "::=" {
				return &ParseError{
					Line:      line,
					Incidence: position.Show(s, pos),
					Err:       eris.Errorf("expected '::=', got '%s...'", s[pos:pos+3]),
				}
			}

			stream.push(&Token{Type: TokenEq})
			pos += 3

		case '|':
			stream.push(&Token{Type: TokenBar})
			pos++

		case '"', '\'':
			// A terminal symbol is enclosed in some form of quotation
			// mark
			val, nextPos, err := until(s, pos+1, s[pos])
			if err != nil {
				return &ParseError{
					Line:      line,
					Incidence: position.Show(s, pos),
					Err:       err,
				}
			}

			stream.push(&Token{Type: TokenT, Value: val})
			pos = nextPos

		case '<':
			// Define a nonterminal symbol as '<...>'
			val, nextPos, err := until(s, pos+1, '>')
			if err != nil {
				return &ParseError{
					Line:      line,
					Incidence: position.Show(s, pos),
					Err:       err,
				}
			}

			stream.push(&Token{Type: TokenNT, Value: val})
			pos = nextPos

		case ' ', '\t', '\r':
			// Whitespace should be skipped
			pos++

		case '\n':
			// Newlines are significant, and have their own token
			stream.push(&Token{Type: TokenEOL})
			pos++

		default:
			return &ParseError{
				Line:      line,
				Incidence: position.Show(s, pos),
				Err:       eris.Errorf("unexpected character '%c'", s[pos]),
			}
		}
	}

	return nil
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

// push will add a token to the token buffer
func (s *TokenStream) push(t *Token) {
	s.buf = append(s.buf, t)
}

// Next will return the next token according to the pos field. It will also
// increment pos to the next place. If we've reached the end of the stream, Next
// will return nil.
func (s *TokenStream) Next() *Token {
	if s.Ended() {
		return nil
	}

	tok := s.buf[s.pos]
	s.pos++

	return tok
}

// Ended will return true if the stream has no more tokens to provide. This is
// effectively the same as testing if `stream.Next() == nil`.
func (s *TokenStream) Ended() bool {
	return s.buf == nil || s.pos >= len(s.buf)
}

// At will return true if the current tokens match the types given. For example,
// `stream.At(TokenNonterminal, TokenOpEqual)` would return true for tokens like
// `<foo> ::=`. We require at least one type, which is why the first parameter
// is not folded into the types slice.
func (s *TokenStream) At(first int, types ...int) bool {
	pos := s.pos
	list := append([]int{first}, types...)

	for _, typ := range list {
		if s.Ended() || s.buf[pos].Type != typ {
			return false
		}

		pos++
	}

	return true
}

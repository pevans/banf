package bnf

import (
	"strings"
	"testing"

	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
)

type errReader int

func (er errReader) Read(_ []byte) (int, error) {
	return 0, eris.New("errReader")
}

var (
	nonterm    = &Token{Type: TokenNT, Value: "aaa"}
	nontermEsc = &Token{Type: TokenNT, Value: "ab>c"}
	term       = &Token{Type: TokenT, Value: "bbb"}
	termEsc    = &Token{Type: TokenT, Value: `d"ef`}
	eq         = &Token{Type: TokenEq}
	bar        = &Token{Type: TokenBar}
	eol        = &Token{Type: TokenEOL}
)

func newStream(types ...int) *TokenStream {
	s := new(TokenStream)

	for _, typ := range types {
		s.push(&Token{Type: typ})
	}

	return s
}

func TestTokenize(t *testing.T) {
	t.Run("failed read", func(t *testing.T) {
		var r errReader

		toks, err := Tokenize(r)
		assert.NotNil(t, toks)
		assert.Error(t, err)
	})

	type tokenTest struct {
		errfn assert.ErrorAssertionFunc
		str   string
		toks  []*Token
	}

	tests := []tokenTest{
		{errfn: assert.NoError, str: `<aaa> ::= "bbb"`, toks: []*Token{nonterm, eq, term}},
		{errfn: assert.NoError, str: `<ab\>c> ::= "d\"ef"`, toks: []*Token{nontermEsc, eq, termEsc}},
		{errfn: assert.NoError, str: `<aaa> ::= "bbb"
<aaa>`, toks: []*Token{nonterm, eq, term, eol, nonterm}},
		{errfn: assert.NoError, str: `<aaa> ::= "bbb" | "d\"ef"`, toks: []*Token{nonterm, eq, term, bar, termEsc}},
		{errfn: assert.NoError, str: `# ignore me pls`, toks: []*Token{}},
		{errfn: assert.NoError, str: `<aaa> ::=#"bbb"`, toks: []*Token{nonterm, eq}},
		{errfn: assert.Error, str: `<aaa> := "bbb"`, toks: []*Token{nonterm}},
		{errfn: assert.Error, str: `<aaa> = what???`, toks: []*Token{nonterm}},
		{errfn: assert.Error, str: `::::==`, toks: []*Token{}},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			toks, err := Tokenize(strings.NewReader(test.str))

			test.errfn(t, err)
			assert.NotNil(t, toks)
			assert.Equal(t, len(test.toks), len(toks.buf))

			for i, tok := range test.toks {
				assert.Equal(t, tok, toks.buf[i])
			}
		})
	}
}

func TestNext(t *testing.T) {
	s := newStream(TokenT, TokenBar)

	// Test that we were able to get the current token, and that pos was updated
	// to the next position
	assert.Equal(t, TokenT, s.Next().Type)
	assert.Equal(t, 1, s.pos)

	// We should see that we iterated to the next token
	assert.Equal(t, TokenBar, s.Next().Type)
	assert.Equal(t, 2, s.pos)

	// There's nothing left, so we should get nil for any further iteration.
	// Additionally, the position shouldn't update.
	assert.Nil(t, s.Next())
	assert.Equal(t, 2, s.pos)
}

func TestEnded(t *testing.T) {
	s := newStream()

	assert.True(t, s.Ended())

	s.push(&Token{Type: TokenEq})
	assert.False(t, s.Ended())

	assert.NotNil(t, s.Next())
	assert.True(t, s.Ended())
}

func TestAt(t *testing.T) {
	s := newStream(TokenNT, TokenEq, TokenT)

	assert.True(t, s.At(TokenNT, TokenEq, TokenT))
	_ = s.Next()
	assert.True(t, s.At(TokenEq, TokenT))

	// This should put us at the end of the stream
	s.pos += 2

	// This should really be false for anything...
	assert.False(t, s.At(TokenEq))
}

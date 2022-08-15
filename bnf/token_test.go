package bnf

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errReader int

func (er errReader) Read(b []byte) (int, error) {
	return 0, errors.New("errReader")
}

var (
	nonterm    = Token{Type: TokenNonterminal, Value: "aaa"}
	nontermEsc = Token{Type: TokenNonterminal, Value: "ab>c"}
	term       = Token{Type: TokenTerminal, Value: "bbb"}
	termEsc    = Token{Type: TokenTerminal, Value: `d"ef`}
	eq         = Token{Type: TokenOpEqual}
	bar        = Token{Type: TokenOpBar}
	eol        = Token{Type: TokenEOL}
)

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
		toks  []Token
	}

	tests := []tokenTest{
		{errfn: assert.NoError, str: `<aaa> ::= "bbb"`, toks: []Token{nonterm, eq, term}},
		{errfn: assert.NoError, str: `<ab\>c> ::= "d\"ef"`, toks: []Token{nontermEsc, eq, termEsc}},
		{errfn: assert.NoError, str: `<aaa> ::= "bbb"
<aaa>`, toks: []Token{nonterm, eq, term, eol, nonterm}},
		{errfn: assert.NoError, str: `<aaa> ::= "bbb" | "d\"ef"`, toks: []Token{nonterm, eq, term, bar, termEsc}},
		{errfn: assert.NoError, str: `# ignore me pls`, toks: []Token{}},
		{errfn: assert.NoError, str: `<aaa> ::=#"bbb"`, toks: []Token{nonterm, eq}},
		{errfn: assert.Error, str: `<aaa> := "bbb"`, toks: []Token{nonterm}},
		{errfn: assert.Error, str: `<aaa> = what???`, toks: []Token{nonterm}},
		{errfn: assert.Error, str: `::::==`, toks: []Token{}},
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

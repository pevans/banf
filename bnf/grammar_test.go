package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const fooBar = `<foo> ::= "bar"`
const fooBarBaz = `
<foo> ::= "bar"
		| "baz"
		`

const manyRules = `
<one> ::=
"1"
<two> ::=

"2"
<three> ::= "3"
`

func TestBuild(t *testing.T) {
	works := func(t *testing.T, bnf string) *Grammar {
		s := new(TokenStream)
		assert.NoError(t, tokenizeString(bnf, s))

		g, err := NewGrammar(s)
		assert.NoError(t, err)
		assert.NotNil(t, g)

		return g
	}

	t.Run("fooBar", func(t *testing.T) {
		g := works(t, fooBar)
		assert.Contains(t, g.Rules, "foo")
	})

	t.Run("fooBarBaz", func(t *testing.T) {
		g := works(t, fooBarBaz)
		assert.Contains(t, g.Rules, "foo")
		assert.NotNil(t, g.Rules["foo"].Condition.OrMatch)
	})

	t.Run("manyRules", func(t *testing.T) {
		g := works(t, manyRules)
		assert.Contains(t, g.Rules, "one")
		assert.Contains(t, g.Rules, "two")
		assert.Contains(t, g.Rules, "three")
	})
}

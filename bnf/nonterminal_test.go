package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNonterminal(t *testing.T) {
	assert.NotNil(t, NewNonterminal(new(Grammar), "abc"))
}

func TestNonterminalMatch(t *testing.T) {
	g, err := NewGrammar(new(TokenStream))
	assert.NoError(t, err)

	s := NewScanner("abc")
	non := NewNonterminal(g, "xyz")

	t.Run("error if no rule found", func(t *testing.T) {
		perr := non.Match(g, s)
		assert.NotNil(t, perr)
		assert.Error(t, perr.Err)
	})

	// missing: test a rule
	rule := NewRule(g, "xyz")
	rule.Condition = new(Expr)
	rule.Condition.Symbols = []Symbol{
		NewTerminal(g, "abc"),
	}

	g.Rules["xyz"] = rule

	t.Run("test a working rule", func(t *testing.T) {
		perr := non.Match(g, s)
		assert.Nil(t, perr)
	})

	rule = NewRule(g, "xyz")
	g.Rules["xyz"] = rule

	t.Run("test a broken rule", func(t *testing.T) {
		perr := non.Match(g, s)
		assert.NotNil(t, perr)
		assert.Error(t, perr.Err)
	})
}

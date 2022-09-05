package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTerminal(t *testing.T) {
	assert.NotNil(t, NewTerminal(new(Grammar), "anything"))
}

func TestTerminalMatch(t *testing.T) {
	s := NewScanner("abc")
	g := new(Grammar)

	term := NewTerminal(g, "ab")

	perr := term.Match(g, s)
	assert.Nil(t, perr)

	t.Run("did it fastforward?", func(t *testing.T) {
		perr := term.Match(g, s)
		assert.NotNil(t, perr)
		assert.NoError(t, perr.Err)

		c := NewTerminal(g, "c")

		perr = c.Match(g, s)
		assert.Nil(t, perr)
	})
}

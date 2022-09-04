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

	did, err := term.Match(g, s)
	assert.NoError(t, err)
	assert.True(t, did)

	t.Run("did it fastforward?", func(t *testing.T) {
		did, err = term.Match(g, s)
		assert.NoError(t, err)
		assert.False(t, did)

		c := NewTerminal(g, "c")

		did, err = c.Match(g, s)
		assert.NoError(t, err)
		assert.True(t, did)
	})
}

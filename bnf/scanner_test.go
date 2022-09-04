package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	s := new(Scanner)

	// We should see that nothing really should match an empty scanner
	assert.False(t, s.StartsWith("haha"))
	assert.True(t, s.StartsWith("")) // except the empty string...

	s.str = "hahaha"
	assert.True(t, s.StartsWith("haha"))

	s.off++
	assert.False(t, s.StartsWith("haha"))
	assert.True(t, s.StartsWith("ahah"))
}

func TestFastForward(t *testing.T) {
	s := new(Scanner)

	// Fast-forwarding an empty scanner should not move the offset
	s.FastForward(111)
	assert.Equal(t, 0, s.off)

	s.str = "haha"

	// We should be able to ffwd to an offset within the bounds of str
	s.FastForward(3)
	assert.Equal(t, 3, s.off)

	// But if we try to go beyond it, we should find ourselves at the end of the
	// string, but no farther
	s.FastForward(3)
	assert.Equal(t, 4, s.off)
}

func TestSave(t *testing.T) {
	s := new(Scanner)

	s.off = 123
	s.Save()

	assert.Equal(t, 123, s.off)
	assert.Equal(t, 123, s.offStack[0])

	s.off++
	s.Save()

	assert.Equal(t, 124, s.off)
	assert.Equal(t, 123, s.offStack[0])
	assert.Equal(t, 124, s.offStack[1])
}

func TestRevert(t *testing.T) {
	s := new(Scanner)

	// If we revert and nothing is in the stack, it should "just work" without
	// causing any issues, even though the offstack length should stay at zero
	assert.NotPanics(t, func() {
		s.Revert()
	})
	assert.Equal(t, 0, len(s.offStack))

	s.Save()
	s.off = 1
	assert.Equal(t, 1, len(s.offStack))

	s.Revert()
	assert.Equal(t, 0, s.off)
	assert.Equal(t, 0, len(s.offStack))
}

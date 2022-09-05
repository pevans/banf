package bnf

import (
	"strings"

	"github.com/pevans/banf/position"
)

// Notes from last night: we should use a string, not a reader
// We might still want to do an offset stack? Unsure. Probably.

type Scanner struct {
	str      string
	off      int
	offStack []int
}

func NewScanner(s string) *Scanner {
	return &Scanner{
		str:      s,
		off:      0,
		offStack: []int{},
	}
}

func (s *Scanner) Show() string {
	return position.Show(s.str, s.off)
}

func (s *Scanner) StartsWith(input string) bool {
	return strings.HasPrefix(s.str[s.off:], input)
}

func (s *Scanner) FastForward(n int) {
	if s.off+n > len(s.str) {
		s.off = len(s.str)
		return
	}

	s.off += n
}

func (s *Scanner) Save() {
	s.offStack = append(s.offStack, s.off)
}

func (s *Scanner) Revert() {
	if len(s.offStack) == 0 {
		return
	}

	end := len(s.offStack) - 1

	s.off = s.offStack[end]
	s.offStack = s.offStack[:end]
}

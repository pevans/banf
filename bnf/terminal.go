package bnf

import "github.com/pevans/banf/blog"

// Terminals are symbols which represent literal string values
type Terminal struct {
	Value string
}

// NewTerminal returns a new Terminal object that has a literal value of val
func NewTerminal(_ *Grammar, val string) *Terminal {
	return &Terminal{
		Value: val,
	}
}

func (t *Terminal) Match(_ *Grammar, scan *Scanner) *ParseError {
	blog.Info("attempting match on '%s'", t.Value)

	match := scan.StartsWith(t.Value)

	if match {
		scan.FastForward(len(t.Value))
	}

	if !match {
		return &ParseError{
			Incidence: scan.Show(),
		}
	}

	return nil
}

package bnf

import "fmt"

// Nonterminals are symbols which represent other rules.
type Nonterminal struct {
	Name string
}

// NewNonterminal returns a new Nonterminal object which is named for val
func NewNonterminal(_ *Grammar, val string) *Nonterminal {
	return &Nonterminal{
		Name: val,
	}
}

func (n *Nonterminal) Match(g *Grammar, scan *Scanner) (bool, error) {
	rule, found := g.Rules[n.Name]
	if !found {
		return false, fmt.Errorf("no such rule <%v>", n.Name)
	}

	match, err := rule.Match(g, scan)
	if err != nil {
		return false, err
	}

	return match, nil
}

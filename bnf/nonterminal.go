package bnf

import "github.com/rotisserie/eris"

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

func (n *Nonterminal) Match(g *Grammar, scan *Scanner) *ParseError {
	rule, found := g.Rules[n.Name]
	if !found {
		return &ParseError{
			Err: eris.Errorf("no such rule <%v>", n.Name),
		}
	}

	perr := rule.Match(g, scan)
	if perr != nil && perr.Err != nil {
		return perr
	}

	if perr != nil {
		return &ParseError{
			Incidence: scan.Show(),
		}
	}

	return nil
}

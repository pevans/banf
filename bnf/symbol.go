package bnf

// Symbols are data that can be matched against input. This interface is
// currently empty because I don't yet know what would be that API.
type Symbol interface {
}

// Terminals are symbols which represent literal string values
type Terminal struct {
	Value string
}

// Nonterminals are symbols which represent other rules.
type Nonterminal struct {
	Name string
}

// NewTerminal returns a new Terminal object that has a literal value of val
func NewTerminal(g *Grammar, val string) *Terminal {
	return &Terminal{
		Value: val,
	}
}

// NewNonterminal returns a new Nonterminal object which is named for val
func NewNonterminal(g *Grammar, val string) *Nonterminal {
	return &Nonterminal{
		Name: val,
	}
}

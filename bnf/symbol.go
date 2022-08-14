package bnf

// A Symbol is one of the components of a BNF rule.
type Symbol struct {
	// Terminal is true when this symbol should be treated literally,
	// and cannot reference another rule.
	Terminal bool

	// Value is the string value of the symbol, regardless of whether
	// it's terminal or not.
	Value string
}

// A Sequence is a set of symbols which comprise one path that can allow
// a Rule to match.
type Sequence []Symbol

func (s *Symbol) Match(g *Grammar, t Token) bool {
	return false
}

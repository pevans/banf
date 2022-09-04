package bnf

// A Rule is a named record that encapsulates some conditional logic for a given
// set of input. It's the entirety of a `<foo> ::= "..."` construct in BNF.
type Rule struct {
	// Name is the name of the rule (the left-hand side of the rule
	// definition)
	Name string

	// Condition is the expression which must be matched for a rule to accept
	// certain input.
	Condition *Expr
}

// NewRule returns a new rule object with an expression already allocated.
// (There is no practical time where you would expect to see a rule without a
// condition.)
func NewRule(_ *Grammar, name string) *Rule {
	r := &Rule{
		Name:      name,
		Condition: new(Expr),
	}

	return r
}

func (r *Rule) Match(g *Grammar, scan *Scanner) (bool, error) {
	return r.Condition.Match(g, scan)
}

package bnf

type Rule struct {
	// Name is the name of the rule (the left-hand side of the rule
	// definition)
	Name string

	// Sequences are the set of Sequences that you can match an input
	// stream and determine if a rule is matched.
	Sequences []Sequence
}

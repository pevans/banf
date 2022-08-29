package bnf

// An Expr is an expression which would match a conditional branch of logic for
// a given input. Every symbol in an expression is evaluated using boolean AND
// logic; boolean ORs by iterating to the next expression in OrMatch.
type Expr struct {
	// Symbols is a slice of Symbols which would be required for us to
	// match.
	Symbols []Symbol

	// OrMatch is an expr which can be considered if this expr is not a
	// match, which itself may link to another expr.
	OrMatch *Expr
}

package bnf

import "github.com/rotisserie/eris"

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

func (e *Expr) Match(g *Grammar, scan *Scanner) *ParseError {
	// You can't really have an expression without any symbols to match, so if
	// it happens, something's wrong
	if len(e.Symbols) == 0 {
		return &ParseError{
			Err: eris.New("expression has no symbols to match"),
		}
	}

	// Save our current position, just in case the match fails
	scan.Save()

	match := true

	for _, sym := range e.Symbols {
		perr := sym.Match(g, scan)
		if perr != nil && perr.Err != nil {
			return perr
		}

		match = match && perr == nil

		// Short-circuit if things didn't work out
		if !match {
			break
		}
	}

	// If we looped over every symbol and they all matched, then we're done and
	// can return here
	if match {
		return nil
	}

	// If we failed to match anything, we can try the next expression
	if e.OrMatch != nil {
		// We should first go back to the position in which we started
		scan.Revert()

		return e.OrMatch.Match(g, scan)
	}

	// At this point, we just need to give up
	return &ParseError{
		Incidence: scan.Show(),
	}
}

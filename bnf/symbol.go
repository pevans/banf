package bnf

// Symbols are data that can be matched against input. This interface is
// currently empty because I don't yet know what would be that API.
type Symbol interface {
	Match(*Grammar, *Scanner) (bool, error)
}

package bnf

type Grammar struct {
	Rules     map[string]Rule
	Terminals []*Symbol
}

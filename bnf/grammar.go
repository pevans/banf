package bnf

import (
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

type Grammar struct {
	MainRule *Rule
	Rules    map[string]*Rule
}

// NewGrammar takes a TokenStream and proceeds to build a grammar from it.
func NewGrammar(stream *TokenStream) (*Grammar, error) {
	g := new(Grammar)
	g.Rules = make(map[string]*Rule)

	if err := g.Build(stream); err != nil {
		return nil, errors.Wrap(err, "couldn't build grammar")
	}

	return g, nil
}

// Build will create all of the rules for a grammar based on an input
// stream of tokens
func (g *Grammar) Build(stream *TokenStream) error {
	var (
		rule *Rule
		expr *Expr
	)

	for {
		token := stream.Next()
		if token == nil {
			break
		}

		switch token.Type {
		case TokenT:
			if expr == nil {
				// At this top level, we don't expect any terminals
				return eris.New("unexpected token: terminal")
			}

			expr.Symbols = append(expr.Symbols, NewTerminal(g, token.Value))

		case TokenNT:
			if expr != nil {
				expr.Symbols = append(expr.Symbols, NewNonterminal(g, token.Value))
				continue
			}

			rule = NewRule(g, token.Value)
			expr = rule.Condition

			if !stream.At(TokenEq) {
				return eris.Errorf("rule must conform to `<%s> ::= ...` syntax", rule.Name)
			}

			// Since we confirmed we're at an OpEqual, let's get rid of it
			_ = stream.Next()

		case TokenEOL:
			if rule != nil && stream.At(TokenNT, TokenEq) {
				// This is the end of the rule definition, and we see another
				// rule being defined
				g.DefineRule(rule)
				rule = nil
				expr = nil
			}

		case TokenBar:
			expr.OrMatch = new(Expr)
			expr = expr.OrMatch

		default:
			return eris.Errorf(
				"unexpected token: %+v", token,
			)
		}
	}

	// We might have reached the end of the string without a newline, but still
	// with a valid rule. We should add it.
	if rule != nil {
		g.DefineRule(rule)
	}

	return nil
}

func (g *Grammar) DefineRule(r *Rule) {
	if g.MainRule == nil {
		g.MainRule = r
	}

	g.Rules[r.Name] = r
}

func (g *Grammar) Match(str string) *ParseError {
	if g.MainRule == nil {
		return &ParseError{
			Err: eris.New("no rules defined for grammar"),
		}
	}

	return g.MainRule.Match(g, NewScanner(str))
}

package bnf

import (
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
)

type ParseError struct {
	File      string
	Line      int
	Incidence string
	Err       error
}

func (p *ParseError) Error() string {
	var b strings.Builder

	if p.File != "" {
		b.WriteString(fmt.Sprintf("%s:\n", p.File))
	}

	if p.Line > 0 {
		b.WriteString(fmt.Sprintf("line %d:\n", p.Line))
	}

	if p.Incidence != "" {
		b.WriteString(p.Incidence + "\n")
	}

	if p.Err != nil {
		b.WriteString(eris.ToString(p.Err, true))
	}

	return b.String()
}

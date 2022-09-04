package bnf

// Terminals are symbols which represent literal string values
type Terminal struct {
	Value string
}

// NewTerminal returns a new Terminal object that has a literal value of val
func NewTerminal(_ *Grammar, val string) *Terminal {
	return &Terminal{
		Value: val,
	}
}

func (t *Terminal) Match(_ *Grammar, scan *Scanner) (bool, error) {
	match := scan.StartsWith(t.Value)

	if match {
		scan.FastForward(len(t.Value))
	}

	return match, nil
}

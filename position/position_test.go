package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShow(t *testing.T) {
	type test struct {
		given  string
		pos    int
		output string
	}

	cases := []test{
		{"this is a test", 5, "this is a test\n     ^"},
	}

	for _, c := range cases {
		assert.Equal(t, c.output, Show(c.given, c.pos))
	}
}

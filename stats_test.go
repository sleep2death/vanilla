package vanilla

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	a := Stats{
		Str: 1,
		Agi: 2,
	}

	b := Stats{
		Str: 1,
		Agi: 2,
		Spi: 3,
	}

	c := a.Add(b)
	assert.Equal(t, 2, c.Str)
	assert.Equal(t, 4, c.Agi)
	assert.Equal(t, 3, c.Spi)
}

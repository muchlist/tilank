package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnum(t *testing.T) {
	m := StApproved
	m2 := StUndefined

	assert.Equal(t, 2, m.EnumIndex())
	assert.Equal(t, "Approved", m.String())
	assert.Equal(t, -1, m2.EnumIndex())
	assert.Equal(t, "Undefined", m2.String())
	assert.Equal(t, StApproved, IntToState(2))
	assert.Equal(t, StUndefined, IntToState(6))
}

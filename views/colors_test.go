package views

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColors_Init(t *testing.T) {
	t.Parallel()
	assert.Len(t, datagramColors, len(datagramColorNames))
}

func TestColors_Markdown(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "fg-blue,bg-orange", markdown("blue", "orange"))
}

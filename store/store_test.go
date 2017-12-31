package store

import (
	"testing"

	"golang.org/x/net/idna"

	"github.com/stretchr/testify/assert"
)

func Test_marshalDomain(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		domain string
		want   string
	}{
		{"Example.orG", "gro.elpmaxe"},
	}
	for _, tt := range tests {
		got, err := marshalDomain(tt.domain, idna.Punycode)
		assert.Equal(tt.want, string(got))
		assert.NoError(err)
	}
}

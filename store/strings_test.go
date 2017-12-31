package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func set(strs ...string) StringSet {
	stringSet := StringSet{}

	for _, str := range strs {
		stringSet.Add(reverse(str))
	}

	return stringSet
}

func Test_stringsFromDomain(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		domain   string
		expected StringSet
	}{
		{"example.org", set("example.org")},
		{"Example.orG", set("example.org")},
		{"example.中国", set("example.xn--fiqs8s")},
		{"ß.example", set("xn--zca.example", "ss.example")},
	}
	for _, tt := range tests {
		got, err := stringsFromDomain(tt.domain)
		assert.Equal(tt.expected, got)
		assert.NoError(err)
	}
}

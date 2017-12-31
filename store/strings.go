package store

import (
	"golang.org/x/net/idna"
)

// StringSet is a simple set of unique strings
type StringSet map[string]struct{}

// Add a string to the set
func (set StringSet) Add(s string) {
	set[s] = struct{}{}
}

// Has returns true if the set contains s
func (set StringSet) Has(s string) bool {
	_, has := set[s]

	return has
}

var profiles = []*idna.Profile{
	idna.Punycode,
	idna.Lookup,
}

// stringsFromDomain returns the various possible encodings of the domain,
// currently the raw punycode and display IDNA profiles. The strings are
// reversed to enable prefix scanning
func stringsFromDomain(domain string) (StringSet, error) {
	set := StringSet{}

	for _, profile := range profiles {
		encoded, err := encodeDomain(domain, profile)
		if err != nil {
			return nil, err
		}
		set.Add(encoded)
	}

	return set, nil
}

func encodeDomain(domain string, profile *idna.Profile) (string, error) {
	ascii, err := profile.ToASCII(domain)
	if err != nil {
		return "", err
	}

	size := len(ascii)
	asciiBytes := make([]byte, size)

	for i := 0; i < size; i++ {
		asciiBytes[size-i-1] = asciiLowerCase(ascii[i])
	}

	return string(asciiBytes), nil
}

func asciiLowerCase(char byte) byte {
	if char >= 'A' && char <= 'Z' {
		return char + 'a' - 'A'
	}

	return char
}

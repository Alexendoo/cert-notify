package store

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/idna"

	"github.com/boltdb/bolt"
	"github.com/gauntface/web-push-go/webpush"
	"golang.org/x/net/publicsuffix"
)

type Store struct {
	db *bolt.DB
}

var (
	// buckets
	bDomains = []byte{0x00}
	bUsers   = []byte{0x01}
)

func New(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bDomains)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bUsers)

		return err
	})

	store := &Store{db}

	return store, err
}

// idna note:
// U-labels seem not to be valid in DNSName fields,
// convert with Punycode + Lookup and if different add two

func (s *Store) AddDomain(d string, sub webpush.Subscription) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		domains := tx.Bucket(bDomains)

		encoded, err := json.Marshal(sub)
		if err != nil {
			return err
		}

		domains.Put([]byte(d), encoded)

		return nil
	})
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DB() *bolt.DB {
	return s.db
}

func Domain(d string) string {
	tld, err := publicsuffix.EffectiveTLDPlusOne(d)
	if err != nil {
		return ""
	}
	return tld
}

func marshalDomain(d string, profile *idna.Profile) ([]byte, error) {
	ascii, err := profile.ToASCII(d)
	b := []byte(strings.ToLower(ascii))

	// reverse b
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return b, err
}

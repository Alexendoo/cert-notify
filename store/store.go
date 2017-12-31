package store

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/gauntface/web-push-go/webpush"
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

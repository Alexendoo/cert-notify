package store

import (
	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

var (
	// buckets
	bDomains = []byte{0x00}
	bUsers   = []byte{0x01}
	bLogs    = []byte{0x02}
)

func New(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		return createBuckets(tx, bDomains, bUsers, bLogs)
	})

	store := &Store{db}

	return store, err
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DB() *bolt.DB {
	return s.db
}

func createBuckets(tx *bolt.Tx, buckets ...[]byte) error {
	for _, bucket := range buckets {
		_, err := tx.CreateBucketIfNotExists(bucket)

		if err != nil {
			return err
		}
	}

	return nil
}

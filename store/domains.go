package store

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gauntface/web-push-go/webpush"
)

// TODO: split to:
// - AddDomains(domains []string, sub *webpush.Subscription)
// - addUser(tx *bolt.Tx, domains _, sub _) <- validate key does not exist here
// - addDomain(tx *bolt.Tx, domain string, sub _)

func (s *Store) AddDomain(d string, sub *webpush.Subscription) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		domains, err := stringsFromDomain(d)
		if err != nil {
			return err
		}

		{ // Add subscription to user bucket
			bUsers := tx.Bucket(bUsers)
			subJSON, err := json.Marshal(sub)
			if err != nil {
				return err
			}
			bUsers.Put(sub.Key, subJSON)
		}

		// Add user references to subscription buckets
		bDomains := tx.Bucket(bDomains)
		for domain := range domains {
			domain := []byte(domain)

			currentBytes := bDomains.Get(domain)

			current := [][]byte{}
			if len(currentBytes) == 0 {
				newBytes, err := json.Marshal([][]byte{sub.Key})
				if err != nil {
					return err
				}

				fmt.Printf("new, %s\n", newBytes)
				return bDomains.Put(domain, newBytes)
			}

			err = json.Unmarshal(currentBytes, &current)
			if err != nil {
				return err
			}

			new := append(current, sub.Key)
			newBytes, err := json.Marshal(new)
			if err != nil {
				return err
			}

			fmt.Printf("append, %s\n", newBytes)
			err = bDomains.Put(domain, newBytes)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

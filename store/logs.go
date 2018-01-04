package store

import (
	"encoding/json"

	"github.com/Alexendoo/cert-notify/ctlog"
	"github.com/boltdb/bolt"
)

func (s *Store) SetLogs(logs []*ctlog.Log) error {
	logsJSON, err := json.Marshal(logs)

	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bLogs).Put(bLogs, logsJSON)
	})
}

func (s *Store) GetLogs() ([]*ctlog.Log, error) {
	logs := []*ctlog.Log{}

	err := s.db.View(func(tx *bolt.Tx) error {
		logsJSON := tx.Bucket(bLogs).Get(bLogs)

		return json.Unmarshal(logsJSON, logs)
	})

	return logs, err
}

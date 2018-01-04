package store

import (
	"encoding/json"

	"github.com/Alexendoo/cert-notify/ctlog"
)

func (s *Store) SetLogs(logs []*ctlog.Log) error {
	logsJSON, err := json.Marshal(logs)

	if err != nil {
		return err
	}

	return s.set(bLogs, bLogs, logsJSON)
}

func (s *Store) GetLogs() ([]*ctlog.Log, error) {
	logs := []*ctlog.Log{}

	logsJSON, err := s.get(bLogs, bLogs)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(logsJSON, logs)

	return logs, err
}

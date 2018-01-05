package ctlog

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/certificate-transparency-go"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
)

// Log is an individual CT log
//
// - https://www.certificate-transparency.org/known-logs
type Log struct {
	// Start index of the next log retrieval, either the end index of the last
	// retrieval or the tree size of a new log
	Index int64
	// DER format of the log's public key used for signature verfication
	Key []byte
	URL string

	c *client.LogClient
	m sync.Mutex
}

func (log *Log) setup(ctx context.Context) (err error) {
	// TODO: not http.DefaultClient
	log.c, err = client.New(log.URL, http.DefaultClient, jsonclient.Options{
		PublicKeyDER: log.Key,
	})
	if err != nil {
		return err
	}
	if log.Index == 0 {
		sth, err := log.c.GetSTH(ctx)

		if err != nil {
			return err
		}

		log.Index = int64(sth.TreeSize)
	}

	return nil
}

type entryFn = func(entry *ct.LogEntry)

func (log *Log) Scan(ctx context.Context, fn entryFn) error {
	log.m.Lock()
	defer log.m.Unlock()

	err := log.setup(ctx)
	if err != nil {
		return err
	}

	delay := 5 * time.Second
	done := ctx.Done()

	for {
		select {
		case <-done:
			return nil
		case <-time.After(delay):
			for log.scan(ctx, fn) {
			}
		}
	}
}

func (log *Log) scan(ctx context.Context, fn entryFn) bool {
	sth, err := log.c.GetSTH(ctx)
	if err != nil {
		return false
	}

	raw, err := log.c.GetRawEntries(ctx, log.Index, int64(sth.TreeSize))
	if err != nil {
		fmt.Println(err, log.Index, int64(sth.TreeSize))
		return false
	}

	done := ctx.Done()
	fmt.Println(len(raw.Entries), log.Index)
	for _, leaf := range raw.Entries {
		// exit if context is cancelled
		select {
		case <-done:
			return false
		default:
		}

		entry, err := ct.LogEntryFromLeaf(log.Index, &leaf)
		if err != nil {
			fmt.Println(log.URL, err)
			continue
		}

		fn(entry)

		log.Index++
	}

	return int64(sth.TreeSize)-log.Index > 100
}

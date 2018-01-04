package ctlog

import (
	"context"
	"fmt"
	"net/http"
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
}

func (log *Log) setupClient() (err error) {
	// TODO: not http.DefaultClient
	log.c, err = client.New(log.URL, http.DefaultClient, jsonclient.Options{
		PublicKeyDER: log.Key,
	})

	return
}

func (log *Log) fetchIndex(ctx context.Context) error {
	sth, err := log.c.GetSTH(ctx)
	if err != nil {
		return err
	}

	log.Index = int64(sth.TreeSize)

	return nil
}

func (log *Log) Scan(ctx context.Context) error {
	if log.c == nil {
		err := log.setupClient()
		if err != nil {
			return err
		}
	}

	if log.Index == 0 {
		err := log.fetchIndex(ctx)
		if err != nil {
			return err
		}
	}

	delay := 5 * time.Second
	done := ctx.Done()

	for {
		select {
		case <-done:
			return nil
		case <-time.After(delay):
			log.scan(ctx)
		}
	}
}

func (log *Log) scan(ctx context.Context) {
	sth, err := log.c.GetSTH(ctx)
	if err != nil {
		return
	}

	raw, err := log.c.GetRawEntries(ctx, log.Index, int64(sth.TreeSize))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, leaf := range raw.Entries {
		index := log.Index + int64(i)
		entry, err := ct.LogEntryFromLeaf(index, &leaf)
		if err != nil {
			fmt.Println(log.URL, err)
			continue
		}

		// TODO: decide on concurrency model for sending notifications to push
		// services

		if entry.X509Cert != nil {
			fmt.Println(entry.X509Cert.DNSNames)
		}
	}
}

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
)

func debase64(encoded string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		panic(err)
	}

	return decoded
}

func main() {
	logger := log.New(os.Stderr, "argon2018 | ", log.Lshortfile)

	ctClient, err := client.New("https://ct.googleapis.com/logs/argon2018", http.DefaultClient, jsonclient.Options{
		Logger:       logger,
		PublicKeyDER: debase64("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0gBVBa3VR7QZu82V+ynXWD14JM3ORp37MtRxTmACJV5ZPtfUA7htQ2hofuigZQs+bnFZkje+qejxoyvk2Q1VaA=="),
	})

	logger.Printf("ctClient: %v\n", ctClient)
	logger.Printf("err: %v\n", err)
}

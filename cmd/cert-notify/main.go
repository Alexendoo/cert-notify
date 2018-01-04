package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Alexendoo/cert-notify/frontend"
	"github.com/Alexendoo/cert-notify/store"
	"github.com/gauntface/web-push-go/webpush"
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
	ctClient, err := client.New("https://ct.googleapis.com/logs/argon2018", http.DefaultClient, jsonclient.Options{
		PublicKeyDER: debase64("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0gBVBa3VR7QZu82V+ynXWD14JM3ORp37MtRxTmACJV5ZPtfUA7htQ2hofuigZQs+bnFZkje+qejxoyvk2Q1VaA=="),
	})

	fmt.Printf("ctClient: %v\n", ctClient)
	fmt.Printf("err: %v\n", err)

	s, err := store.New("z.db")
	fmt.Println(err)
	sub := &webpush.Subscription{
		Auth:     []byte{0x30, 0x32, 0x53, 0x9F, 0xe1, 0xff, 0x00},
		Endpoint: "https://example.org/push/awioefioawioef",
		Key:      []byte("Secret key"),
	}
	fmt.Println(s.AddDomain("macleod.io", sub))

	frontend.Serve()
}

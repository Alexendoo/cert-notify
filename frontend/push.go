package frontend

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gauntface/web-push-go/webpush"
)

var curve = elliptic.P256()

func subscribe(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sub, err := webpush.SubscriptionFromJSON(bytes)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("sub: %#+v\n", sub)
}

func getPublicKey(w http.ResponseWriter, r *http.Request) {
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	public := private.PublicKey
	pubBytes := elliptic.Marshal(curve, public.X, public.Y)

	w.Write(pubBytes)
}

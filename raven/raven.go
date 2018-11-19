package raven

import (
	"crypto/tls"

	"github.com/ravendb/ravendb-go-client"
	"log"
	"os"
)

var Store *ravendb.DocumentStore

func SetupStore() {
	Store = ravendb.NewDocumentStoreWithUrlAndDatabase("https://a.vvanm.ravendb.community/", "")
	Store.SetDatabase("gymtracker")

	if os.Getenv("PORT") != "" {
		crt := os.Getenv("RAVEN_CERT")
		key := os.Getenv("RAVEN_KEY")

		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			log.Println(err)
		}
		Store.SetCertificate(&ravendb.KeyStore{
			Certificates: []tls.Certificate{cert},
		})

	} else {
		cert, err := tls.LoadX509KeyPair("raven/vvanm.crt", "raven/vvanm.key")
		if err != nil {
			log.Println(err)
		}

		Store.SetCertificate(&ravendb.KeyStore{
			Certificates: []tls.Certificate{cert},
		})

	}

	Store.Initialize()
}

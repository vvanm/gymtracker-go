package raven

import (
	"crypto/tls"
	"github.com/ravendb/ravendb-go-client"
	"log"
)

var Store *ravendb.DocumentStore

func SetupStore() {
	Store = ravendb.NewDocumentStoreWithUrlAndDatabase("https://a.vvanm.ravendb.community/", "")
	Store.SetDatabase("gymtracker")

	//if os.Getenv("PORT") != "" {

	cert, err := tls.LoadX509KeyPair("raven/vvanm.crt", "raven/vvanm.key")
	if err != nil {
		log.Println(err)
	}

	Store.SetCertificate(&ravendb.KeyStore{
		Certificates: []tls.Certificate{cert},
	})

	//	}

	Store.Initialize()
}

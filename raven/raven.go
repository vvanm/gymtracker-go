package raven

import (
	"crypto/tls"
	"fmt"
	"github.com/ravendb/ravendb-go-client"
	"log"
	"os"
	"strings"
)

var Store *ravendb.DocumentStore

func SetupStore() {
	Store = ravendb.NewDocumentStoreWithUrlAndDatabase("https://a.vvanm.ravendb.community/", "")
	Store.SetDatabase("gymtracker")

	// fetcha all env variables
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}

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

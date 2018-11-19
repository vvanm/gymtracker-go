package raven

import (
	"github.com/ravendb/ravendb-go-client"
)

var Store *ravendb.DocumentStore

func SetupStore() {
	Store = ravendb.NewDocumentStoreWithUrlAndDatabase("http://localhost:8080", "")
	Store.SetDatabase("gymtracker")
	Store.Initialize()
}

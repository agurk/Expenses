package main

import (
	"b2/backend"
	"b2/components/analysis"
	"b2/components/changes"
	"b2/components/exrecords"
	"b2/components/managed/accounts"
	"b2/components/managed/classifications"
	"b2/components/managed/docexmappings"
	"b2/components/managed/documents"
	"b2/components/managed/expenses"
	"b2/components/suggestions"
	"b2/manager"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Config is to hold the data loaded at runtime from the config file
type Config struct {
	Host         string
	ServerCert   string
	ServerKey    string
	DB           string
	SwUser       uint64
	SwToken      string
	DocsLocation string
}

func loadConfig() *Config {
	config := new(Config)
	file, err := os.Open("config")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}

type handler interface {
	Handle(w http.ResponseWriter, req *http.Request)
	Path() string
	LongPath() string
}

func addHandler(h handler) {
	http.HandleFunc(h.Path(), h.Handle)
	http.HandleFunc(h.LongPath(), h.Handle)
}

func main() {
	config := loadConfig()

	backend := backend.Instance(config.DB)
	backend.Accounts = accounts.Instance(backend)
	backend.Classifications = classifications.Instance(backend)
	backend.Documents = documents.Instance(backend)
	backend.Expenses = expenses.Instance(backend)
	backend.Mappings = docexmappings.Instance(backend)
	backend.Splitwise.BearerToken = config.SwToken
	backend.Splitwise.User = config.SwUser
	backend.DocsLocation = config.DocsLocation
	backend.Start()

	addHandler(analysis.Instance("/analysis", backend.DB))
	addHandler(manager.Instance("/documents", backend.Documents))
	addHandler(manager.Instance("/expenses", backend.Expenses))
	addHandler(manager.Instance("/expenses/accounts", backend.Accounts))
	addHandler(manager.Instance("/expenses/classifications", backend.Classifications))
	addHandler(exrecords.Instance("/expenses/externalrecords", backend))
	addHandler(suggestions.Instance("/expenses/suggestions", backend))
	addHandler(manager.Instance("/mappings", backend.Mappings))

	http.HandleFunc("/processor", backend.Process)

	c := changes.Instance("/changes/", backend)
	http.HandleFunc(c.Path, c.Handle)

	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	log.Fatal(http.ListenAndServeTLS(config.Host, config.ServerCert, config.ServerKey, nil))
}

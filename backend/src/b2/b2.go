package main

import (
	"b2/backend"
	"b2/components/analysis"
	"b2/components/exrecords"
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

type Config struct {
	Host       string
	ServerCert string
	ServerKey  string
	DB         string
	SW_User    uint64
	SW_Token   string
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
	GetPath() string
	GetLongPath() string
}

func addHandler(h handler) {
	http.HandleFunc(h.GetPath(), h.Handle)
	http.HandleFunc(h.GetLongPath(), h.Handle)
}

func main() {
	config := loadConfig()

	backend := backend.Instance(config.DB)
	backend.Documents = documents.Instance(backend)
	backend.Expenses = expenses.Instance(backend)
	backend.Classifications = classifications.Instance(backend)
	backend.Mappings = docexmappings.Instance(backend)
	backend.Splitwise.User = config.SW_User
	backend.Splitwise.BearerToken = config.SW_Token

	addHandler(analysis.Instance("/analysis", backend.DB))
	addHandler(manager.Instance("/documents", backend.Documents))
	addHandler(manager.Instance("/expenses", backend.Expenses))
	addHandler(manager.Instance("/expenses/classifications", backend.Classifications))
	addHandler(exrecords.Instance("/expenses/externalrecords", backend))
	addHandler(suggestions.Instance("/expenses/suggestions", backend))
	addHandler(manager.Instance("/mappings", backend.Mappings))

	http.HandleFunc("/processor", backend.Process)

	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	log.Fatal(http.ListenAndServeTLS(config.Host, config.ServerCert, config.ServerKey, nil))
}

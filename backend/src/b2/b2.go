package main

import (
	"b2/backend"
	"b2/components/analysis"
	"b2/components/exrecords"
	"b2/components/managed/classifications"
	"b2/components/managed/docexmappings"
	"b2/components/managed/documents"
	"b2/components/managed/expenses"
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

func main() {
	config := loadConfig()

	backend := backend.Instance(config.DB)
	backend.Documents = documents.Instance(backend)
	backend.Expenses = expenses.Instance(backend)
	backend.Classifications = classifications.Instance(backend)
	backend.Mappings = docexmappings.Instance(backend)
	backend.Splitwise.User = config.SW_User
	backend.Splitwise.BearerToken = config.SW_Token

	docWebManager := new(manager.WebHandler)
	docWebManager.Initalize("/documents/", backend.Documents)

	exWebManager := new(manager.WebHandler)
	exWebManager.Initalize("/expenses/", backend.Expenses)

	clWebManager := new(manager.WebHandler)
	clWebManager.Initalize("/expense_classifications/", backend.Classifications)

	mapWebManager := new(manager.WebHandler)
	mapWebManager.Initalize("/mappings/", backend.Mappings)

	analWebManager := new(analysis.WebHandler)
	analWebManager.Initalize(backend.DB)

	recManager := new(exrecords.WebHandler)
	recManager.Initalize(backend)

	http.HandleFunc("/expenses/classifications", clWebManager.MultipleHandler)
	http.HandleFunc("/expenses/classifications/", clWebManager.IndividualHandler)
	http.HandleFunc("/expenses/externalrecords/", recManager.Handler)
	http.HandleFunc("/expenses/", exWebManager.IndividualHandler)
	http.HandleFunc("/expenses", exWebManager.MultipleHandler)
	http.HandleFunc("/documents/", docWebManager.IndividualHandler)
	http.HandleFunc("/documents", docWebManager.MultipleHandler)
	http.HandleFunc("/mappings/", mapWebManager.IndividualHandler)

	http.HandleFunc("/analysis/", analWebManager.Handler)
	http.HandleFunc("/processor", backend.Process)

	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	log.Fatal(http.ListenAndServeTLS(config.Host, config.ServerCert, config.ServerKey, nil))
}

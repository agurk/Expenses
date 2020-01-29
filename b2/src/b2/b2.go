package main

import (
	"b2/backend"
	"b2/components/analysis"
	"b2/components/managed/classifications"
	"b2/components/managed/docexmappings"
	"b2/components/managed/documents"
	"b2/components/managed/expenses"
	"b2/manager"
	"log"
	"net/http"
)

func main() {
	backend := backend.Instance("/home/timothy/src/Expenses/expenses.db")
	backend.Documents = documents.Instance(backend)
	backend.Expenses = expenses.Instance(backend)
	backend.Classifications = classifications.Instance(backend)
	backend.Mappings = docexmappings.Instance(backend)

	docWebManager := new(manager.WebHandler)
	exWebManager := new(manager.WebHandler)
	clWebManager := new(manager.WebHandler)
	mapWebManager := new(manager.WebHandler)
	analWebManager := new(analysis.WebHandler)

	analWebManager.Initalize(backend.DB)
	docWebManager.Initalize("/documents/", backend.Documents)
	exWebManager.Initalize("/expenses/", backend.Expenses)
	mapWebManager.Initalize("/mappings/", backend.Mappings)
	clWebManager.Initalize("/expense_classifications/", backend.Classifications)

	http.HandleFunc("/expenses/classifications", clWebManager.MultipleHandler)
	http.HandleFunc("/expenses/classifications/", clWebManager.IndividualHandler)
	http.HandleFunc("/expenses/", exWebManager.IndividualHandler)
	http.HandleFunc("/expenses", exWebManager.MultipleHandler)
	http.HandleFunc("/documents/", docWebManager.IndividualHandler)
	http.HandleFunc("/documents", docWebManager.MultipleHandler)
	http.HandleFunc("/mappings/", mapWebManager.IndividualHandler)

	http.HandleFunc("/analysis/", analWebManager.Handler)
	http.HandleFunc("/processor", backend.Process)

	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	log.Fatal(http.ListenAndServeTLS("localhost:8000", "certs/server.crt", "certs/server.key", nil))
}

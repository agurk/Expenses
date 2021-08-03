package documents

import (
	"b2/backend"
	"b2/components/changes"
	"b2/components/managed/docexmappings"
	"b2/components/managed/expenses"
	"b2/errors"
	"b2/manager"
	"b2/moneyutils"
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/schema"
)

// Query represents the search criteria that can be used when looking for a document
type Query struct {
	// both of these are toggling only that value
	Starred   bool `schema:"starred"`
	Unmatched bool `schema:"unmatched"`
	Archived  bool `schema:"archived"`
}

// DocManager is a component used by a manager to manager documents
type DocManager struct {
	backend *backend.Backend
}

// Instance returns an instantiated caching manager configured for documents
func Instance(backend *backend.Backend) manager.Manager {
	dm := new(DocManager)
	dm.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(dm)
	return general
}

// Load returns a document that matches the passed in id, if extant
func (dm *DocManager) Load(did uint64) (manager.Thing, error) {
	return loadDocument(did, dm.backend.DB)
}

// AfterLoad adds the mappings to expenses. This will replace/reload any
// that have already been loaded so this can be called when there have been
// changes to the mappings
func (dm *DocManager) AfterLoad(thing manager.Thing) error {
	document := Cast(thing)
	v := new(docexmappings.Query)
	v.DocumentID = document.ID
	mapps, err := dm.backend.Mappings.Find(v)
	document.Lock()
	defer document.Unlock()
	document.Expenses = []*docexmappings.Mapping{}
	for _, thing := range mapps {
		mapping := docexmappings.Cast(thing)
		document.Expenses = append(document.Expenses, mapping)
	}
	return errors.Wrap(err, "documents.AfterLoad")
}

// Find returns a slice of ids for all documents that match the criteria in the Query
// or a url.Values encoded version of it
func (dm *DocManager) Find(query interface{}) ([]uint64, error) {
	var search *Query
	switch query.(type) {
	case *Query:
		search = query.(*Query)
	case url.Values:
		search = new(Query)
		decoder := schema.NewDecoder()
		err := decoder.Decode(search, query.(url.Values))
		if err != nil {
			return nil, errors.Wrap(err, "documents.Find")
		}
	default:
		panic("Unknown type passed to find function")
	}
	return findDocuments(search, dm.backend.DB)
}

// FindExisting does nothing for documents
func (dm *DocManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create saves a new version of the document into the db
func (dm *DocManager) Create(thing manager.Thing) error {
	document := Cast(thing)
	err := createDocument(document, dm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "documents.Create")
	}
	// todo: this seems like an inefficient way to get the document processed
	dm.backend.ReprocessDocument <- document.ID
	dm.backend.Change <- changes.DocumentEvent
	return nil
}

func (dm *DocManager) ocr(doc *Document) error {
	cmd := exec.Command("tesseract", dm.backend.DocsLocation+"/"+doc.Filename, "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "documents.ocr")
	}
	doc.Lock()
	doc.Text = fmt.Sprintf("%s", out.String())
	doc.Unlock()
	return nil
}

func (dm *DocManager) matchExpenses(doc *Document) error {
	doc.Lock()
	defer doc.Unlock()
	dates := make(map[string]bool)
	year := "(2?0?[0-9]{2})"
	month := "(0?[0-9]|1?[0-2])"
	day := "([12][0-9]|3[01]|0?[0-9])"
	seperators := `[-–—\/\\. ]`
	date := regexp.MustCompile(year + seperators + month + seperators + day)
	for _, values := range date.FindAllStringSubmatch(doc.Text, -2) {
		dates[makeDateString(values[1], values[2], values[3])] = true
	}
	date = regexp.MustCompile(day + seperators + month + seperators + year)
	for _, values := range date.FindAllStringSubmatch(doc.Text, -2) {
		dates[makeDateString(values[3], values[2], values[1])] = true
	}
	date = regexp.MustCompile(month + seperators + day + seperators + year)
	for _, values := range date.FindAllStringSubmatch(doc.Text, -2) {
		dates[makeDateString(values[3], values[1], values[2])] = true
	}

	if len(dates) == 0 {
		return nil
	}

	query := new(expenses.Query)
	for key := range dates {
		query.Dates = append(query.Dates, key)
	}
	exes, err := dm.backend.Expenses.Find(query)
	if err != nil {
		return errors.Wrap(err, "documents.matchExpenses")
	}
	results := make([]uint64, len(exes))
	var wg sync.WaitGroup
	for i, ex := range exes {
		wg.Add(1)
		go func(thing manager.Thing, results []uint64, pos int) {
			defer wg.Done()
			expense := expenses.Cast(thing)
			expense.RLock()
			defer expense.RUnlock()
			for _, term := range strings.Split(expense.Description, " ") {
				// also skips ""
				if len(term) < 2 {
					continue
				}
				if strings.Contains(strings.ToLower(doc.Text), strings.ToLower(term)) {
					results[pos]++
				}
			}
			var amountStr string
			if expense.FX.Amount != 0 {
				// todo: decimial places?
				amountStr = fmt.Sprintf("%f", expense.FX.Amount)
			} else {
				amountStr, err = moneyutils.StringAbs(expense.Amount, expense.Currency)
				if err != nil {
					errors.Print(errors.Wrap(err, "expenses.matchExpenses"))
				}
			}
			// this should work as the decimial seperator is . so will do a wildcard search
			match, _ := regexp.MatchString(amountStr, doc.Text)
			if match {
				results[pos]++
			}

		}(ex, results, i)
	}
	wg.Wait()
	var maxVal uint64
	for _, val := range results {
		if val > maxVal {
			maxVal = val
		}
	}
	if maxVal == 0 {
		return nil
	}
	for i, val := range results {
		if val == maxVal {
			mapping := new(docexmappings.Mapping)
			mapping.EID = exes[i].GetID()
			mapping.DID = doc.ID
			err := dm.backend.Mappings.New(mapping)
			if err != nil {
				errors.Print(err)
			}
			// the document will have its mappings updated after this by calling
			// the After load function again
			dm.backend.ReloadExpenseMappings <- mapping.EID
		}
	}
	return nil
}

func makeDateString(year, month, day string) string {
	if len(year) == 2 {
		// todo: fix before the year 2100
		year = "20" + year
	}
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return year + "-" + month + "-" + day
}

// Update causes the db to be updated with any changes to the document
func (dm *DocManager) Update(thing manager.Thing) error {
	document := Cast(thing)
	dm.backend.Change <- changes.DocumentEvent
	return updateDocument(document, dm.backend.DB)
}

// NewThing returns a newly instantiated empty unsaved document
func (dm *DocManager) NewThing() manager.Thing {
	return new(Document)
}

// Combine is not implemented for documents
func (dm *DocManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "documents.Combine", true)
}

// Delete removes the document from the db
func (dm *DocManager) Delete(thing manager.Thing) error {
	document := Cast(thing)
	document.Lock()
	defer document.Unlock()
	err := deleteDocument(document, dm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "documents.Delete")
	}
	document.deleted = true
	for _, expense := range document.Expenses {
		err = dm.backend.Mappings.Delete(expense)
		if err != nil {
			errors.Print(err)
		}
	}
	dm.backend.Change <- changes.DocumentEvent
	return errors.Wrap(err, "documents.Delete")
}

// Process will reperform OCR on a document and reclassify it
func (dm *DocManager) Process(id uint64) {
	doc, err := dm.backend.Documents.Get(id)
	if err != nil {
		errors.Print(err)
		return
	}
	document := Cast(doc)
	if document.Text == "" {
		err = dm.ocr(document)
		if err != nil {
			errors.Print(err)
			return
		}
		err = dm.Update(document)
		if err != nil {
			errors.Print(err)
			return
		}
	}
	err = dm.matchExpenses(document)
	if err != nil {
		errors.Print(err)
		return
	}
	dm.AfterLoad(document)
	dm.backend.Change <- changes.DocumentEvent
}

// ReclassifyAll will reclassify all documents that have not got confirmed matches or are archived
func (dm *DocManager) ReclassifyAll() error {
	eligible, err := reclassifyableDocs(dm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "documents.ReclassifyAll")
	}
	for _, id := range eligible {
		d, err := dm.Load(id)
		if err != nil {
			return errors.Wrap(err, "documents.ReclassifyAll")
		}
		doc := Cast(d)
		err = dm.matchExpenses(doc)
		if err != nil {
			return errors.Wrap(err, "documents.ReclassifyAll")
		}
	}
	dm.backend.Change <- changes.DocumentEvent
	return nil
}

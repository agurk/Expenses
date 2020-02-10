package documents

import (
	"b2/backend"
	"b2/components/changes"
	"b2/components/managed/docexmappings"
	"b2/components/managed/expenses"
	"b2/errors"
	"b2/manager"
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/schema"
)

type Query struct {
	// both of these are toggling only that value
	Starred   bool `schema:"starred"`
	Unmatched bool `schema:"unmatched"`
	Archived  bool `schema:"archived"`
}

type DocManager struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) manager.Manager {
	dm := new(DocManager)
	dm.initalize(backend)
	general := new(manager.CachingManager)
	general.Initalize(dm)
	return general
}

func (dm *DocManager) initalize(backend *backend.Backend) {
	dm.backend = backend
}

func (dm *DocManager) Load(did uint64) (manager.Thing, error) {
	return loadDocument(did, dm.backend.DB)
}

func (dm *DocManager) AfterLoad(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		panic("Non document passed to function")
	}
	v := new(docexmappings.Query)
	v.DocumentID = document.ID
	mapps, err := dm.backend.Mappings.Find(v)
	document.Lock()
	defer document.Unlock()
	document.Expenses = []*docexmappings.Mapping{}
	for _, thing := range mapps {
		mapping, ok := thing.(*(docexmappings.Mapping))
		if !ok {
			panic("Non mapping returned from function")
		}
		document.Expenses = append(document.Expenses, mapping)
	}
	return errors.Wrap(err, "documents.AfterLoad")
}

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

func (dm *DocManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (dm *DocManager) Create(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		panic("Non document passed to function")
	}
	err := createDocument(document, dm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "documents.Create")
	}
	dm.backend.DocumentsProcessChan <- document.ID
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

	fmt.Println("dates:", dates)
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
		go func(expens manager.Thing, results []uint64, pos int) {
			defer wg.Done()
			expense, ok := expens.(*expenses.Expense)
			if !ok {
				panic("Non expense sent to function")
			}
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
			var amount float64
			if expense.FX.Amount != 0 {
				amount = expense.FX.Amount
			} else {
				amount = float64(expense.Amount) / 100
			}
			if strings.Contains(fmt.Sprintf("%f", amount), doc.Text) {
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
			dm.backend.ExpensesDepsChan <- mapping.EID
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

func (dm *DocManager) Update(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		panic("Non document passed to function")
	}
	dm.backend.Change <- changes.DocumentEvent
	return updateDocument(document, dm.backend.DB)
}

func (dm *DocManager) NewThing() manager.Thing {
	return new(Document)
}

func (dm *DocManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "documents.Combine")
}

func (dm *DocManager) Delete(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		panic("Non document passed to function")
	}
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

func (dm *DocManager) Process(id uint64) {
	doc, err := dm.backend.Documents.Get(id)
	document, ok := doc.(*Document)
	if !ok {
		panic("Non document passed to function")
		return
	}
	if err != nil {
		errors.Print(err)
		return
	}
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

func (dm *DocManager) ReclassifyAll() error {
	eligible, err := getReclassifyableDocs(dm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "documents.ReclassifyAll")
	}
	fmt.Println(eligible)
	for _, id := range eligible {
		d, err := dm.Load(id)
		if err != nil {
			return errors.Wrap(err, "documents.ReclassifyAll")
		}
		doc, ok := d.(*Document)
		if !ok {
			panic("Not document passed to function")
		}
		err = dm.matchExpenses(doc)
		if err != nil {
			return errors.Wrap(err, "documents.ReclassifyAll")
		}
	}
	dm.backend.Change <- changes.DocumentEvent
	return nil
}

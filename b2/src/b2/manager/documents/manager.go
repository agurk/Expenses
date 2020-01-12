package documents

import (
	"b2/backend"
	"b2/manager"
	"b2/manager/docexmappings"
	"b2/manager/expenses"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

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
		return errors.New("Non document passed to function")
	}
	v := new(docexmappings.Query)
	v.DocumentId = document.ID
	mapps, err := dm.backend.Mappings.Find(v)
	for _, thing := range mapps {
		mapping, ok := thing.(*(docexmappings.Mapping))
		if !ok {
			return errors.New("Non mapping returned from function")
		}
		document.Expenses = append(document.Expenses, mapping)
	}
	return err
}

func (dm *DocManager) Find(params interface{}) ([]uint64, error) {
	return findDocuments(dm.backend.DB)
}

func (dm *DocManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (dm *DocManager) Create(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		return errors.New("Non document passed to function")
	}
	err := createDocument(document, dm.backend.DB)
	if err != nil {
		return err
	}
	return dm.matchExpenses(document)
}

func (dm *DocManager) matchExpenses(doc *Document) error {
	doc.RLock()
	defer doc.RUnlock()
	dates := make(map[string]bool)
	year := "(2?0?[0-9]{2})"
	month := "(0?[0-9]|1?[0-2])"
	day := "([12][0-9]|3[01]|0?[0-9])"
	seperators := `[-–—\/\\.]`
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

	query := new(expenses.Query)
	for key, _ := range dates {
		query.Dates = append(query.Dates, key)
	}
	exes, err := dm.backend.Expenses.Find(query)
	if err != nil {
		return nil
	}
	results := make([]uint64, len(exes))
	var wg sync.WaitGroup
	for i, ex := range exes {
		wg.Add(1)
		go func(expens manager.Thing, results []uint64, pos int) {
			defer wg.Done()
			expense, ok := expens.(*expenses.Expense)
			if !ok {
				fmt.Println("Non expense sent to function")
				return
			}
			expense.RLock()
			defer expense.RUnlock()
			for _, term := range strings.Split(expense.Description, " ") {
				if term == "" {
					continue
				}
				if strings.Contains(strings.ToLower(doc.Text), strings.ToLower(term)) {
					results[pos]++
				}
			}
			if strings.Contains(fmt.Sprintf("%f", expense.Amount), doc.Text) {
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
	for i, val := range results {
		if val == maxVal {
			mapping := new(docexmappings.Mapping)
			mapping.EID = exes[i].GetID()
			mapping.DID = doc.ID
			err := dm.backend.Mappings.Save(mapping)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(doc.ID, exes[i].GetID())
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
		return errors.New("Non document passed to function")
	}
	return updateDocument(document, dm.backend.DB)
}

func (dm *DocManager) NewThing() manager.Thing {
	return new(Document)
}

func (dm *DocManager) Combine(one, two manager.Thing) error {
	return errors.New("Not implemented")
}

func (dm *DocManager) Delete(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		return errors.New("Non document passed to function")
	}
	document.Lock()
	defer document.Unlock()
	err := deleteDocument(document, dm.backend.DB)
	if err != nil {
		return err
	}
	document.deleted = true
	for _, expense := range document.Expenses {
		// todo: err getting masked
		err = dm.backend.Mappings.Delete(expense)
	}
	return err
}

package expenses

import (
	"b2/backend"
	"b2/manager"
	"b2/manager/docexmappings"
	"errors"
	"github.com/gorilla/schema"
	"math"
	"net/url"
	"regexp"
	"strings"
)

type Query struct {
	From string `schema:"from"`
	To   string `schema:"to"`
	// Date can be completed, but will not be used directly, instead to & from
	// will take its value
	Date           string   `schema:"date"`
	Search         string   `schema:"search"`
	Dates          []string `schema:"dates"`
	Classification string   `schema:"classification"`
}

func cleanQuery(query *Query) {
	if query.Date != "" {
		query.From = query.Date
		query.To = query.Date
	}
	classRE := regexp.MustCompile(`classification:"([^"]*)"`)
	for _, value := range classRE.FindAllStringSubmatch(query.Search, -2) {
		query.Classification = value[1]
	}
}

type ExManager struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) manager.Manager {
	em := new(ExManager)
	em.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(em)
	return general
}

func (em *ExManager) Load(eid uint64) (manager.Thing, error) {
	return loadExpense(eid, em.backend.DB)
}

func (em *ExManager) AfterLoad(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	v := new(docexmappings.Query)
	v.ExpenseId = expense.ID
	mapps, err := em.backend.Mappings.Find(v)
	for _, thing := range mapps {
		mapping, ok := thing.(*(docexmappings.Mapping))
		if !ok {
			return errors.New("Non mapping returned from function")
		}
		expense.Documents = append(expense.Documents, mapping)
	}
	return err
}

func (em *ExManager) Find(query interface{}) ([]uint64, error) {
	var search *Query
	switch query.(type) {
	case *Query:
		search = query.(*Query)
	case url.Values:
		search = new(Query)
		decoder := schema.NewDecoder()
		err := decoder.Decode(search, query.(url.Values))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown type passed to find function")
	}
	cleanQuery(search)
	if search.Classification != "" {
		return findExpensesClassification(search, em.backend.DB)
	}
	if search.Search != "" {
		return findExpensesSearch(search, em.backend.DB)
	}
	if len(search.Dates) > 0 {
		return findExpensesDates(search, em.backend.DB)
	}
	return findExpensesDate(search, em.backend.DB)
}

func (em *ExManager) FindExisting(thing manager.Thing) (uint64, error) {
	expense, ok := thing.(*Expense)
	if !ok {
		return 0, errors.New("Non expense passed to function")
	}
	expense.RLock()
	defer expense.RUnlock()
	if expense.TransactionReference != "" {
		oldEid, err := findExpenseByTranRef(expense.TransactionReference, expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, err
		} else if oldEid > 0 {
			return oldEid, nil
		}
	}
	if expense.Metadata.Temporary {
		oldEid, err := findExpenseByDetails(expense.Amount, expense.Date, expense.Description, expense.Currency, expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, err
		} else if oldEid > 0 {
			return oldEid, nil
		}
	} else {
		// todo: improve matching (date range? tipping percent? ignore description spaces?)
		results, err := getTempExpenseDetails(expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, err
		}
		lastDiff := 10000000.0
		confirmedTolerance := 0.05
		var eid uint64 = 0
		for _, result := range results {
			// check same sign
			if expense.Amount*result.Amount < 0 {
				continue
			}
			diff := math.Abs(math.Abs(result.Amount)-math.Abs(expense.Amount)) / math.Abs(expense.Amount)
			if diff > confirmedTolerance {
				continue
			}
			oldDesc := strings.ToLower(strings.Replace(expense.Description, " ", "", -1))
			newDesc := strings.ToLower(strings.Replace(result.Description, " ", "", -1))
			if oldDesc != newDesc {
				continue
			}
			if diff < lastDiff {
				eid = result.ID
				lastDiff = diff
			}
		}
		return eid, nil
	}
	return 0, nil
}

func (em *ExManager) Create(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	em.classifyExpense(expense)
	return createExpense(expense, em.backend.DB)
}

func (em *ExManager) classifyExpense(expense *Expense) {
	// todo: add some logic here
	expense.Metadata.Classification = 5
	expense.Metadata.Confirmed = false
}

func (em *ExManager) Combine(ex, ex2 manager.Thing) error {
	expense, ok := ex.(*Expense)
	exMergeWith, ok2 := ex2.(*Expense)
	if !(ok && ok2) {
		return errors.New("Non expense passed to function")
	}
	expense.Merge(exMergeWith)
	exMergeWith.deleted = true
	for _, mapping := range exMergeWith.Documents {
		mapping.EID = expense.ID
		// todo: deal with error?
		em.backend.Mappings.Save(mapping)
	}
	exMergeWith.Documents = nil
	expense.Documents = nil
	return em.AfterLoad(expense)
}

func (em *ExManager) Update(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	return updateExpense(expense, em.backend.DB)
}

func (em *ExManager) Delete(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	expense.Lock()
	defer expense.Unlock()
	err := deleteExpense(expense, em.backend.DB)
	if err != nil {
		return nil
	}
	expense.deleted = true
	for _, mapping := range expense.Documents {
		// todo err getting masked
		err = em.backend.Mappings.Delete(mapping)
	}
	return err
}

func (em *ExManager) NewThing() manager.Thing {
	return new(Expense)
}

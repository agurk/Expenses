package expenses

import (
	"b2/manager"
	"b2/manager/docexmappings"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	From string
	To   string
}

type ExManager struct {
	db          *sql.DB
	docMappings manager.Manager
}

func Instance(db *sql.DB, docMappings manager.Manager) manager.Manager {
	em := new(ExManager)
	em.db = db
	em.docMappings = docMappings
	general := new(manager.CachingManager)
	general.Initalize(em)
	return general
}

func (em *ExManager) Load(eid uint64) (manager.Thing, error) {
	return loadExpense(eid, em.db)
}

func (em *ExManager) AfterLoad(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	v := url.Values{}
	v.Set("expense", strconv.FormatUint(expense.ID, 10))
	mapps, err := em.docMappings.GetMultiple(v)
	for _, thing := range mapps {
		mapping, ok := thing.(*(docexmappings.Mapping))
		if !ok {
			return errors.New("Non mapping returned from function")
		}
		expense.Documents = append(expense.Documents, mapping)
	}
	if err != nil {
		return err
	}
	return err
}

func (em *ExManager) FindFromUrl(params url.Values) ([]uint64, error) {
	var query Query
	for key, elem := range params {
		fmt.Println(key)
		fmt.Println(elem)
		// Query() returns empty string as value when no value set for key
		if len(elem) != 1 || elem[0] == "" {
			return nil, errors.New("Invalid query parameter " + key)
		}
		switch key {
		case "date":
			// todo: validate date
			query.From = elem[0]
			query.To = elem[0]
		case "from":
			query.From = elem[0]
		case "to":
			query.To = elem[0]
		default:
			return nil, errors.New("Invalid query parameter " + key)
		}
	}

	if query.To == "" || query.From == "" {
		return nil, errors.New("Missing date in date range")
	}
	return em.Find(&query)
}

func (em *ExManager) Find(query *Query) ([]uint64, error) {
	return findExpenses(query, em.db)
}

func (em *ExManager) FindExisting(thing manager.Thing) (uint64, error) {
	expense, ok := thing.(*Expense)
	if !ok {
		return 0, errors.New("Non expense passed to function")
	}
	expense.RLock()
	defer expense.RUnlock()
	if expense.TransactionReference != "" {
		oldEid, err := findExpenseByTranRef(expense.TransactionReference, expense.AccountID, em.db)
		if err != nil {
			return 0, err
		} else if oldEid > 0 {
			return oldEid, nil
		}
	}
	if expense.Metadata.Temporary {
		oldEid, err := findExpenseByDetails(expense.Amount, expense.Date, expense.Description, expense.Currency, expense.AccountID, em.db)
		if err != nil {
			return 0, err
		} else if oldEid > 0 {
			return oldEid, nil
		}
	} else {
		// todo: improve matching (date range? tipping percent? ignore description spaces?)
		results, err := getTempExpenseDetails(expense.AccountID, em.db)
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
	err := createExpense(expense, em.db)
	if err != nil {
		return err
	}
	return nil
}

func (em *ExManager) classifyExpense(expense *Expense) {
	// todo: add some logic here
	expense.Metadata.Classification = 19
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
		em.docMappings.Save(mapping)
	}
	return nil
}

func (em *ExManager) Update(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	return updateExpense(expense, em.db)
}

func (em *ExManager) Delete(ex manager.Thing) error {
	expense, ok := ex.(*Expense)
	if !ok {
		return errors.New("Non expense passed to function")
	}
	// todo: delete document mappings
	return deleteExpense(expense, em.db)
}

func (em *ExManager) NewThing() manager.Thing {
	return new(Expense)
}

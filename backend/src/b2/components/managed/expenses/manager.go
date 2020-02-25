package expenses

import (
	"b2/backend"
	"b2/components/changes"
	"b2/components/managed/docexmappings"
	"b2/errors"
	"b2/manager"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/schema"
)

// Query contains all the paramaters that can be used to search for an expense
type Query struct {
	From            string   `schema:"from"`
	To              string   `schema:"to"`
	Date            string   `schema:"date"`
	Search          string   `schema:"search"`
	Dates           []string `schema:"dates"`
	Classification  string   `schema:"classification"`
	OnlyUnconfirmed bool     `schema:"unconfirmed"`
	OnlyTemporary   bool     `schema:"temporary"`
}

func findQueryParams(query *Query) {
	classRE := regexp.MustCompile(`clas[sifcaton]{0,10}: *(?:"([^"]*)"|([^ ]*))`)
	value := classRE.FindStringSubmatch(query.Search)
	if len(value) >= 3 {
		// add them together as either the first or second match should be empty
		query.Classification = value[1] + value[2]
		query.Search = classRE.ReplaceAllString(query.Search, "$1")
	}
	conRE := regexp.MustCompile("(not: *con[firmed]{0,6})")
	value = conRE.FindStringSubmatch(query.Search)
	if len(value) >= 1 {
		query.OnlyUnconfirmed = true
		query.Search = conRE.ReplaceAllString(query.Search, "")
	}
	tempRE := regexp.MustCompile("(is: *tem[poray]{0,6})")
	value = tempRE.FindStringSubmatch(query.Search)
	if len(value) >= 1 {
		query.OnlyTemporary = true
		query.Search = tempRE.ReplaceAllString(query.Search, "")
	}
}

// ExManager is the component used by a manager to manage expenses
type ExManager struct {
	backend *backend.Backend
}

// Instance returns an instantiated caching manager configured for expenses
func Instance(backend *backend.Backend) manager.Manager {
	em := new(ExManager)
	em.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(em)
	return general
}

// Load returns an expense (if extant) for a specific ID
func (em *ExManager) Load(eid uint64) (manager.Thing, error) {
	expense, err := loadExpense(eid, em.backend.DB)
	if err != nil {
		return nil, errors.Wrap(err, "expenses.Load")
	}
	if expense.Metadata.Classification == 0 {
		classifyExpense(expense, em.backend.DB)
		err = em.Update(expense)
	}
	return expense, errors.Wrap(err, "expenses.Load")
}

// AfterLoad performs the loading of dependencies to the expense, like mappings to
// documents. This function will replace any previous loaded values so can be called again
// if they need to be reloaded
func (em *ExManager) AfterLoad(thing manager.Thing) error {
	expense := Cast(thing)
	v := new(docexmappings.Query)
	v.ExpenseID = expense.ID
	mapps, err := em.backend.Mappings.Find(v)
	expense.Lock()
	defer expense.Unlock()
	expense.Documents = []*docexmappings.Mapping{}
	for _, thing := range mapps {
		mapping := docexmappings.Cast(thing)
		expense.Documents = append(expense.Documents, mapping)
	}
	return errors.Wrap(err, "expenses.AfterLoad")
}

// Find returns a list of IDs relating to either a Query object or the same parameters
// encoded in a url.Values
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
			return nil, errors.Wrap(err, "expenses.Find")
		}
	default:
		return nil, errors.New("Unexpected type passed to function", nil, "expenses.Find", false)
	}
	findQueryParams(search)
	return findExpenses(search, em.backend.DB)
}

// FindExisting will return an expenses that match the one passed in, used to avoid saving duplicates
// Temporary expenses will be matched more greedily than confirmed ones
func (em *ExManager) FindExisting(thing manager.Thing) (uint64, error) {
	var oldEid uint64 = 0
	var err error
	expense := Cast(thing)
	expense.RLock()
	defer expense.RUnlock()
	if expense.TransactionReference != "" {
		oldEid, err = findExpenseByTranRef(expense.TransactionReference, expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.FindExisting")
		}
	} else if expense.Metadata.Temporary {
		oldEid, err = findExpenseByDetails(expense.Amount, expense.Date, expense.Description, expense.Currency, expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.FindExisting")
		}
	}
	if oldEid == 0 {
		// todo: improve matching (date range? tipping percent? ignore description spaces?)
		results, err := tempExpenseDetails(expense.AccountID, em.backend.DB)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.FindExisting")
		}
		lastDiff := 10000000.0
		confirmedTolerance := 0.05
		for _, result := range results {
			// check same sign
			if expense.Amount*result.Amount < 0 {
				continue
			}
			resAmt := float64(result.Amount)
			exAmt := float64(expense.Amount)
			diff := math.Abs(math.Abs(resAmt)-math.Abs(exAmt)) / math.Abs(exAmt)
			if diff > confirmedTolerance {
				continue
			}
			currentDesc := strings.ToLower(strings.Replace(expense.Description, " ", "", -1))
			currentDescRe := regexp.MustCompile(regexp.QuoteMeta(currentDesc))
			newDesc := strings.ToLower(strings.Replace(result.Description, " ", "", -1))
			newDescRe := regexp.MustCompile(regexp.QuoteMeta(newDesc))
			// todo: allow somewhat matching strings?
			if !currentDescRe.MatchString(newDesc) && !newDescRe.MatchString(currentDesc) {
				continue
			}
			if diff < lastDiff {
				oldEid = result.ID
				lastDiff = diff
			}
		}
	}
	// Logic for what to return is to make sure only a temporary expense is overwritten
	// and a duplicate expense is met with an error
	// | NewEx | OldEx | Return |
	// --------------------------
	// |  T    |  T    | EID    | Updating Temp
	// |  P    |  T    | EID    | Updating Temp to Permanent
	// |  T    |  P    | Err    | New Temp for Duplicate
	// |  P    |  P    | Err    | Duplicate
	if oldEid > 0 {
		oldEx, err := em.Load(oldEid)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.FindExisting")
		}
		// if this can't be cast to an expense, something has gone very wrong
		if oldEx.(*Expense).Metadata.Temporary {
			return oldEid, nil
		} else if expense.Metadata.Temporary {
			return 0, errors.New(fmt.Sprintf("Could not create new temporary expense, as expense already exists %d", oldEid), errors.Forbidden, "expenses.FindExisting", true)
		} else {
			return 0, errors.New(fmt.Sprintf("Could not create new expense, as expense already exists %d", oldEid), errors.Forbidden, "expenses.FindExisting", true)
		}
	}
	return 0, nil
}

// Create saves new version of the passed in expense in the db
func (em *ExManager) Create(thing manager.Thing) error {
	expense := Cast(thing)
	classifyExpense(expense, em.backend.DB)
	err := createExpense(expense, em.backend.DB)
	if err != nil {
		return errors.Wrap(err, "expenses.Create")
	}
	em.backend.Change <- changes.ExpenseEvent
	return nil
}

// Combine merges the two expenses either normally or as a commission depending on what
// behaviour is specified in the params
func (em *ExManager) Combine(thing, thing2 manager.Thing, params string) error {
	expense := Cast(thing)
	exMergeWith := Cast(thing2)
	if params == "commission" {
		expense.MergeAsCommission(exMergeWith)
	} else {
		expense.Merge(exMergeWith)
	}
	exMergeWith.deleted = true
	for _, mapping := range exMergeWith.Documents {
		mapping.EID = expense.ID
		err := em.backend.Mappings.Save(mapping)
		if err != nil {
			errors.Print(err)
		}
	}
	exMergeWith.Documents = nil
	expense.Documents = nil
	em.backend.Change <- changes.ExpenseEvent
	return em.AfterLoad(expense)
}

// Update saves any changes made to the expense into the db
// and alerts the backend of a change being made
func (em *ExManager) Update(thing manager.Thing) error {
	expense := Cast(thing)
	err := updateExpense(expense, em.backend.DB)
	em.backend.Change <- changes.ExpenseEvent
	return errors.Wrap(err, "expenses.Update")
}

// Delete the expense from the DB
func (em *ExManager) Delete(thing manager.Thing) error {
	expense := Cast(thing)
	expense.Lock()
	defer expense.Unlock()
	err := deleteExpense(expense, em.backend.DB)
	if err != nil {
		return nil
	}
	expense.deleted = true
	for _, mapping := range expense.Documents {
		err = em.backend.Mappings.Delete(mapping)
		if err != nil {
			errors.Print(err)
		}
	}
	em.backend.Change <- changes.ExpenseEvent
	return errors.Wrap(err, "expenses.Delete")
}

// NewThing returns a newly instatiated empty unsave expense
func (em *ExManager) NewThing() manager.Thing {
	return new(Expense)
}

// Process will reclassify the expense
func (em *ExManager) Process(id uint64) {
	ex, err := em.backend.Expenses.Get(id)
	if err != nil {
		errors.Print(err)
		return
	}
	expense := Cast(ex)
	classifyExpense(expense, em.backend.DB)
	err = em.Update(expense)
	if err != nil {
		fmt.Println("Error updating expense")
	}
	em.backend.Change <- changes.ExpenseEvent
}

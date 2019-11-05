package expenses 

import "database/sql"
import "fmt"
import "errors"

type dbExpense struct {
    ID uint64
    TransactionReference sql.NullString
    Description sql.NullString
    DetailedDescription sql.NullString
    AccountID int
    Date sql.NullString
    ProcessDate sql.NullString
    Amount float64
    Currency string
    Commission  sql.NullInt64
    MetaModified sql.NullString
    MetaTemp sql.NullBool
    MetaConfirmed sql.NullBool
    MetaClassi sql.NullString
    FXAmnt sql.NullFloat64
    FXCCY sql.NullString
    FXRate sql.NullFloat64
}

type dbClassification struct {
    ID uint64
    Description string
    Hidden bool
    From string
    To sql.NullString
}

func parseSQLstr(str *sql.NullString) string {
        if !str.Valid {
            return ""
        }
        return str.String
}

func parseSQLint(integer *sql.NullInt64) int64 {
        if !integer.Valid {
            return 0
        }
        return integer.Int64
}

func parseSQLfloat(flt *sql.NullFloat64) float64 {
        if !flt.Valid {
            return 0
        }
        return flt.Float64
}

func parseSQLbool(boolean *sql.NullBool) bool {
        if !boolean.Valid {
            return false 
        }
        return boolean.Bool
}

func cleanDate(date string) string {
    // horrible hack
    if date == "" {
        return date
    }
    return date[0:len("1234-12-12")]
}

func result2expense(result *dbExpense) *Expense {
    expense := new(Expense)
    // mandatory fields 
    expense.ID = result.ID
    expense.AccountID = result.AccountID
    expense.Amount = result.Amount
    expense.Currency = result.Currency
    // Optional fields
    expense.TransactionReference = parseSQLstr(&result.TransactionReference)
    expense.Description = parseSQLstr(&result.Description)
    expense.DetailedDescription = parseSQLstr(&result.DetailedDescription)
    expense.Commission = parseSQLint(&result.Commission)
    expense.Date = cleanDate(parseSQLstr(&result.Date))
    expense.ProcessDate = cleanDate(parseSQLstr(&result.ProcessDate))
    expense.FX.Amount = parseSQLfloat(&result.FXAmnt)
    expense.FX.Currency = parseSQLstr(&result.FXCCY)
    expense.FX.Rate = parseSQLfloat(&result.FXRate)
    expense.Metadata.Confirmed = parseSQLbool(&result.MetaConfirmed)
    //expense.Metadata.Tagged = parseSQL
    expense.Metadata.Temporary = parseSQLbool(&result.MetaTemp)
    expense.Metadata.Modified = parseSQLstr(&result.MetaModified)
    expense.Metadata.Classification = parseSQLstr(&result.MetaClassi)
    return expense
}

func result2classification(result *dbClassification) *Classification {
    classification := new(Classification)
    classification.ID = result.ID
    classification.Description = result.Description
    classification.From = result.From
    classification.To = parseSQLstr(&result.To)
    classification.Hidden = result.Hidden
    return classification
}

func findExpenses(from, to string, db *sql.DB) ([]uint64, error) {
    rows, err := db.Query("select eid from expenses where date between $1 and $2", from, to)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var eids []uint64
    for rows.Next() {
        var eid uint64
        err = rows.Scan(&eid)
        if err != nil {
            return nil, err
        }
        eids = append(eids, eid)
    }
    return eids, err
}

func getClassifications(db *sql.DB) ([]*Classification, error) {
    rows, err := db.Query("select cid, name, validfrom, validto, isexpense from classificationdef")
    if err != nil {
        return nil, err
    }
    var classifications []*Classification
    defer rows.Close()
    for rows.Next() {
        class := new(dbClassification)
        err = rows.Scan(&class.ID,
                        &class.Description,
                        &class.From,
                        &class.To,
                        &class.Hidden)
        if err != nil {
            return nil, err
        }
        classifications = append(classifications, result2classification(class))

    }
    return classifications, err
}

func loadExpense(eid uint64, db *sql.DB) (*Expense, error) {
    rows, err := db.Query("select e.aid, e.description, e.amount, e.ccy, e.amountFX, e.ccyFX, e.fxRate, e.commission, e.date, e.modified, e.temporary, e.reference, e.detaileddescription, c.cid, c.confirmed, e.processDate from expenses e, classifications c where e.eid = $1 and e.eid = c.eid", eid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    expense := new(dbExpense)
    if rows.Next() {
        err = rows.Scan(&expense.AccountID,
                        &expense.Description,
                        &expense.Amount,
                        &expense.Currency,
                        &expense.FXAmnt,
                        &expense.FXCCY,
                        &expense.FXRate,
                        &expense.Commission,
                        &expense.Date,
                        &expense.MetaModified,
                        &expense.MetaTemp,
                        &expense.TransactionReference,
                        &expense.DetailedDescription,
                        &expense.MetaClassi,
                        &expense.MetaConfirmed,
                        &expense.ProcessDate)
        expense.ID = eid
    } else {
        return nil, errors.New("404")
    }
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return result2expense(expense), nil
}

func createExpense(e *Expense, db *sql.DB) error {
    // todo: check values are legit before writing
    _, err := db.Exec("insert into expenses (aid, description, amount, ccy, amountFX, ccyFX, fxRate, commission, date, temporary, reference, detaileddescription, processDate) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", e.AccountID, e.Description, e.Amount, e.Currency, e.FX.Amount, e.FX.Currency, e.FX.Rate, e.Commission, e.Date, e.Metadata.Temporary, e.TransactionReference, e.DetailedDescription, e.ProcessDate)
    // todo: make this safer!
    rows, err := db.Query("select max(eid) from expenses")
    if err != nil {
        return err
    }
    if rows.Next() {
        var eid uint64
        err = rows.Scan(&eid)
        if err != nil {
            return err
        }
        e.ID = eid
        fmt.Println(eid)
        rows.Close()
    } else {
        return errors.New("Error creating new expense")
    }

    _, err = db.Exec("delete from classifications where eid = $1; insert into classifications  (eid, cid, confirmed) values ($2, $3, $4)", e.ID, e.ID, e.Metadata.Classification, e.Metadata.Confirmed)


    return err
}

func updateExpenes(e *Expense, db *sql.DB) error {
    e.RLock()
    defer e.RUnlock()
    // Todo: Check values are legit before writing
    _, err := db.Exec("update expenses set aid = $1, description = $2, amount = $3, ccy = $4, amountFX = $5, ccyFX = $6, fxRate = $7, commission = $8, date = $9, temporary = $10, reference = $11, detaileddescription = $12, processDate = $13 where eid = $14; delete from classifications where eid = $15; insert into classifications  (eid, cid, confirmed) values ($16, $17, $18)", e.AccountID, e.Description, e.Amount, e.Currency, e.FX.Amount, e.FX.Currency, e.FX.Rate, e.Commission, e.Date, e.Metadata.Temporary, e.TransactionReference, e.DetailedDescription, e.ProcessDate, e.ID, e.ID, e.ID, e.Metadata.Classification, e.Metadata.Confirmed)
    if err != nil {
        return err
    }
    return err
}

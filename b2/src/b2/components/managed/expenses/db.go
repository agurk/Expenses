package expenses

import (
	"database/sql"
	"errors"
	"fmt"
)

type dbExpense struct {
	ID                   uint64
	TransactionReference sql.NullString
	Description          sql.NullString
	DetailedDescription  sql.NullString
	AccountID            uint
	Date                 sql.NullString
	ProcessDate          sql.NullString
	Amount               int64
	Currency             string
	Commission           int64
	MetaModified         sql.NullString
	MetaTemp             sql.NullBool
	MetaConfirmed        sql.NullBool
	MetaClassi           sql.NullInt64
	MetaOldValues        sql.NullString
	FXAmnt               sql.NullFloat64
	FXCCY                sql.NullString
	FXRate               sql.NullFloat64
}

type expenseDetails struct {
	ID          uint64
	Description string
	Amount      int64
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
	// todo improve date handling
	if date == "" {
		return date
	}
	return date[0:10]
}

func result2expense(result *dbExpense) *Expense {
	expense := new(Expense)
	// mandatory fields
	expense.ID = result.ID
	expense.AccountID = result.AccountID
	expense.Amount = result.Amount
	expense.Currency = result.Currency
	expense.Commission = result.Commission
	// Optional fields
	expense.TransactionReference = parseSQLstr(&result.TransactionReference)
	expense.Description = parseSQLstr(&result.Description)
	expense.DetailedDescription = parseSQLstr(&result.DetailedDescription)
	expense.Date = cleanDate(parseSQLstr(&result.Date))
	expense.ProcessDate = cleanDate(parseSQLstr(&result.ProcessDate))
	expense.FX.Amount = parseSQLfloat(&result.FXAmnt)
	expense.FX.Currency = parseSQLstr(&result.FXCCY)
	expense.FX.Rate = parseSQLfloat(&result.FXRate)
	expense.Metadata.Confirmed = parseSQLbool(&result.MetaConfirmed)
	//expense.Metadata.Tagged = parseSQL
	expense.Metadata.Temporary = parseSQLbool(&result.MetaTemp)
	expense.Metadata.Modified = parseSQLstr(&result.MetaModified)
	expense.Metadata.Classification = parseSQLint(&result.MetaClassi)
	expense.Metadata.OldValues = parseSQLstr(&result.MetaOldValues)
	return expense
}

func findExpenses(query *Query, db *sql.DB) ([]uint64, error) {
	var args []interface{}
	dbQuery := `
		select
			e.eid
		from
			expenses e,
			classifications c,
			classificationdef cd
		where
			c.cid = cd.cid
			and c.eid = e.eid`
	if query.Search != "" {
		args = append(args, "%"+query.Search+"%")
		dbQuery += fmt.Sprintf(" and description like $%d", len(args))
	}
	if query.Classification != "" {
		args = append(args, "%"+query.Classification+"%")
		dbQuery += fmt.Sprintf(" and cd.name like $%d", len(args))
	}
	if query.Date != "" {
		args = append(args, query.Date)
		dbQuery += fmt.Sprintf(" and date = $%d", len(args))
	}
	if query.From != "" {
		args = append(args, query.From)
		dbQuery += fmt.Sprintf(" and date >= $%d", len(args))
	}
	if query.To != "" {
		args = append(args, query.To)
		dbQuery += fmt.Sprintf(" and date <= $%d", len(args))
	}
	if len(query.Dates) > 0 {
		var instr string
		for _, date := range query.Dates {
			if instr != "" {
				instr += ","
			}
			args = append(args, date)
			instr += fmt.Sprintf("$%d", len(args))
		}
		dbQuery += `and date in(` + instr + ")"
	}
	rows, err := db.Query(dbQuery, args...)
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
	return eids, nil
}

func findExpenseByTranRef(ref string, account uint, db *sql.DB) (uint64, error) {
	rows, err := db.Query("select eid from expenses where Reference = $1 and aid = $2", ref, account)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var eid uint64
	// todo : what about results with multiple tran refs?
	for rows.Next() {
		err = rows.Scan(&eid)
		if err != nil {
			return 0, err
		}
	}
	return eid, err
}

func findExpenseByDetails(amount int64, date, description, currency string, account uint, db *sql.DB) (uint64, error) {
	rows, err := db.Query(`
		select
			eid
		from
			expenses
		where
			aid = $1
			and date = $2
			and description = $3
            and amount = $4
			and ccy = $5`,
		account, date, description, amount, currency)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var eid uint64
	// todo : what about results with multiple results
	for rows.Next() {
		err = rows.Scan(&eid)
		if err != nil {
			return 0, err
		}
	}
	return eid, err
}

func getTempExpenseDetails(account uint, db *sql.DB) ([]*expenseDetails, error) {
	rows, err := db.Query(`
		select
			eid,
			amount,
			description
		from
			expenses
		where
			aid = $1
			and temporary
		order by
			date asc`, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	temprows := []*expenseDetails{}
	for rows.Next() {
		row := new(expenseDetails)
		err = rows.Scan(&row.ID, &row.Amount, &row.Description)
		if err != nil {
			return nil, err
		}
		temprows = append(temprows, row)
	}
	return temprows, err
}

func loadExpense(eid uint64, db *sql.DB) (*Expense, error) {
	rows, err := db.Query(`
        select
            e.aid,
            e.description,
            e.amount,
            e.ccy,
            e.amountFX,
            e.ccyFX,
            e.fxRate,
            e.commission,
            e.date,
            e.modified,
            e.temporary,
            e.reference,
            e.detaileddescription,
            c.cid,
            c.confirmed,
            e.processDate,
			e.oldValues
        from
            expenses e
        left join
            classifications c
			on e.eid = c.eid
        where
            e.eid = $1`,
		eid)
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
			&expense.ProcessDate,
			&expense.MetaOldValues)
		expense.ID = eid
	} else {
		return nil, errors.New("404")
	}
	if err != nil {
		return nil, err
	}
	fullExpense := result2expense(expense)
	err = addExternalRecords(fullExpense, db)
	return fullExpense, err
}

func loadDocuments(e *Expense, db *sql.DB) ([]uint64, error) {
	rows, err := db.Query(`
        select
            dem.did
        from
            DocumentExpenseMapping dem
        where
            dem.eid = $1`,
		e.ID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	dids := []uint64{}
	for rows.Next() {
		var did uint64
		err = rows.Scan(&did)
		if err != nil {
			return nil, err
		}
		dids = append(dids, did)

	}
	return dids, err
}

func createExpense(e *Expense, db *sql.DB) error {
	err := e.Check()
	if err != nil {
		return err
	}
	e.Lock()
	defer e.Unlock()
	res, err := db.Exec(`insert into
							expenses (
								aid,
								description,
								amount,
								ccy,
								amountFX,
								ccyFX,
								fxRate,
								commission,
								date,
								temporary,
								reference,
								detaileddescription,
								processDate,
								oldValues)
							values
								($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		e.AccountID,
		e.Description,
		e.Amount,
		e.Currency,
		e.FX.Amount,
		e.FX.Currency,
		e.FX.Rate,
		e.Commission,
		e.Date,
		e.Metadata.Temporary,
		e.TransactionReference,
		e.DetailedDescription,
		e.ProcessDate,
		e.Metadata.OldValues)

	if err != nil {
		return err
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		e.ID = uint64(rid)
	} else {
		return errors.New("Error creating new expense")
	}

	_, err = db.Exec("delete from classifications where eid = $1; insert into classifications  (eid, cid, confirmed) values ($2, $3, $4)", e.ID, e.ID, e.Metadata.Classification, e.Metadata.Confirmed)
	if err != nil {
		return err
	}
	return saveExternalRecords(e, db)
}

func updateExpense(e *Expense, db *sql.DB) error {
	err := e.Check()
	if err != nil {
		return err
	}
	e.RLock()
	defer e.RUnlock()
	_, err = db.Exec(`
		update
			expenses
		set
			aid = $1,
			description = $2,
            amount = $3,
            ccy = $4,
            amountFX = $5,
            ccyFX = $6,
            fxRate = $7,
            commission = $8,
            date = $9,
            temporary = $10,
            reference = $11,
            detaileddescription = $12,
            processDate = $13,
            oldValues = $14
		where
			eid = $15;

		delete from
			classifications
		where
			eid = $16;

		insert into
			classifications
				(eid, cid, confirmed)
			values
				($17, $18, $19)`,
		e.AccountID,
		e.Description,
		e.Amount,
		e.Currency,
		e.FX.Amount,
		e.FX.Currency,
		e.FX.Rate,
		e.Commission,
		e.Date,
		e.Metadata.Temporary,
		e.TransactionReference,
		e.DetailedDescription,
		e.ProcessDate,
		e.Metadata.OldValues,
		e.ID,
		e.ID,
		e.ID,
		e.Metadata.Classification,
		e.Metadata.Confirmed)
	if err != nil {
		return err
	}
	return saveExternalRecords(e, db)
}

func saveExternalRecords(e *Expense, db *sql.DB) error {
	// assuming that the expense we're given is already locked
	_, err := db.Exec(`
		delete from
			ExternalRecords
		where
			eid = $1`,
		e.ID)
	if err != nil {
		return err
	}
	for _, ref := range e.ExternalRecords {
		_, err = db.Exec(`
			insert into
				ExternalRecords (eid, Type, Reference, FullAmount)
			values
				($1, $2, $3, $4)`,
			e.ID, ref.Type, ref.Reference, ref.FullAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

func addExternalRecords(e *Expense, db *sql.DB) error {
	e.ExternalRecords = []*ExternalRecord{}
	rows, err := db.Query(`
		select
			type,
			reference,
			fullamount
		from
			ExternalRecords
		where
			eid = $1`,
		e.ID)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var typeValue, reference string
		var oldamount int64
		err = rows.Scan(&typeValue, &reference, &oldamount)
		if err != nil {
			return err
		}
		extRec := new(ExternalRecord)
		extRec.Type = typeValue
		extRec.Reference = reference
		extRec.FullAmount = oldamount
		e.ExternalRecords = append(e.ExternalRecords, extRec)
	}
	return nil
}

func deleteExpense(e *Expense, db *sql.DB) error {
	// assuming that the expense we're given is already locked
	_, err := db.Exec(`
		delete from expenses where eid = $1;
		delete from classifications where eid = $2;
		delete from externalrecords where eid = $3`,
		e.ID, e.ID, e.ID)
	return err
}

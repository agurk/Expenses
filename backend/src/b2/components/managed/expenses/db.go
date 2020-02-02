package expenses

import (
	"b2/errors"
	"database/sql"
	"fmt"
)

type expenseDetails struct {
	ID          uint64
	Description string
	Amount      int64
}

func cleanDate(date string) string {
	// todo improve date handling
	if date == "" {
		return date
	}
	return date[0:10]
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
	if query.OnlyUnconfirmed {
		dbQuery += " and not c.confirmed "
	}
	if query.OnlyTemporary {
		dbQuery += " and e.temporary "
	}
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
		dbQuery += ` and date in(` + instr + ")"
	}
	rows, err := db.Query(dbQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "expenses.findExpenses")
	}
	defer rows.Close()
	var eids []uint64
	for rows.Next() {
		var eid uint64
		err = rows.Scan(&eid)
		if err != nil {
			return nil, errors.Wrap(err, "expenses.findExpenses")
		}
		eids = append(eids, eid)
	}
	return eids, nil
}

func findExpenseByTranRef(ref string, account uint, db *sql.DB) (uint64, error) {
	rows, err := db.Query("select eid from expenses where Reference = $1 and aid = $2", ref, account)
	if err != nil {
		return 0, errors.Wrap(err, "expenses.findExpenseByTranRef")
	}
	defer rows.Close()
	var eid uint64
	// todo : what about results with multiple tran refs?
	for rows.Next() {
		err = rows.Scan(&eid)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.findExpenseByTranRef")
		}
	}
	return eid, errors.Wrap(err, "expenses.findExpenseByTranRef")
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
		return 0, errors.Wrap(err, "expenses.findExpenseByDetails")
	}
	defer rows.Close()
	var eid uint64
	// todo : what about results with multiple results
	for rows.Next() {
		err = rows.Scan(&eid)
		if err != nil {
			return 0, errors.Wrap(err, "expenses.findExpenseByDetails")
		}
	}
	return eid, errors.Wrap(err, "expenses.findExpenseByDetails")
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
		return nil, errors.Wrap(err, "expenses.getTempExpenseDetails")
	}
	defer rows.Close()
	temprows := []*expenseDetails{}
	for rows.Next() {
		row := new(expenseDetails)
		err = rows.Scan(&row.ID, &row.Amount, &row.Description)
		if err != nil {
			return nil, errors.Wrap(err, "expenses.getTempExpenseDetails")
		}
		temprows = append(temprows, row)
	}
	return temprows, errors.Wrap(err, "expenses.getTempExpenseDetails")
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
		return nil, errors.Wrap(err, "expenses.loadExpenses")
	}
	defer rows.Close()
	//expense := new(dbExpense)
	expense := new(Expense)
	if rows.Next() {
		err = rows.Scan(&expense.AccountID,
			&expense.Description,
			&expense.Amount,
			&expense.Currency,
			&expense.FX.Amount,
			&expense.FX.Currency,
			&expense.FX.Rate,
			&expense.Commission,
			&expense.Date,
			&expense.Metadata.Modified,
			&expense.Metadata.Temporary,
			&expense.TransactionReference,
			&expense.DetailedDescription,
			&expense.Metadata.Classification,
			&expense.Metadata.Confirmed,
			&expense.ProcessDate,
			&expense.Metadata.OldValues)
		expense.ID = eid
	} else {
		return nil, errors.New(fmt.Sprintf("Expense %d not found", eid), errors.ThingNotFound, "expenses.loadExpenses")
	}
	if err != nil {
		return nil, errors.Wrap(err, "expenses.loadExpenses")
	}
	expense.Date = cleanDate(expense.Date)
	expense.ProcessDate = cleanDate(expense.ProcessDate)
	err = addExternalRecords(expense, db)
	return expense, errors.Wrap(err, "expenses.loadExpenses")
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
		return nil, errors.Wrap(err, "expenses.loadDocuments")
	}
	dids := []uint64{}
	for rows.Next() {
		var did uint64
		err = rows.Scan(&did)
		if err != nil {
			return nil, errors.Wrap(err, "expenses.loadDocuments")
		}
		dids = append(dids, did)

	}
	return dids, errors.Wrap(err, "expenses.loadDocument")
}

func createExpense(e *Expense, db *sql.DB) error {
	err := e.Check()
	if err != nil {
		return errors.Wrap(err, "expenses.createExpense")
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
								oldValues,
								modified)
							values
								($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
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
		e.Metadata.Modified)

	if err != nil {
		return errors.Wrap(err, "expenses.createExpense")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		e.ID = uint64(rid)
	} else {
		return errors.New("Error creating new expense", errors.InternalError, "expenses.createExpense")
	}

	_, err = db.Exec("delete from classifications where eid = $1; insert into classifications  (eid, cid, confirmed) values ($2, $3, $4)", e.ID, e.ID, e.Metadata.Classification, e.Metadata.Confirmed)
	if err != nil {
		return errors.Wrap(err, "expenses.createExpense")
	}
	return saveExternalRecords(e, db)
}

func updateExpense(e *Expense, db *sql.DB) error {
	err := e.Check()
	if err != nil {
		return errors.Wrap(err, "expenses.updateExpense")
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
            oldValues = $14,
			modified = $15
		where
			eid = $16;

		delete from
			classifications
		where
			eid = $17;

		insert into
			classifications
				(eid, cid, confirmed)
			values
				($18, $19, $20)`,
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
		e.Metadata.Modified,
		e.ID,
		e.ID,
		e.ID,
		e.Metadata.Classification,
		e.Metadata.Confirmed)
	if err != nil {
		return errors.Wrap(err, "expenses.updateExpense")
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
		return errors.Wrap(err, "expenses.saveExternalRecords")
	}
	for _, ref := range e.ExternalRecords {
		_, err = db.Exec(`
			insert into
				ExternalRecords (eid, Type, Reference, FullAmount)
			values
				($1, $2, $3, $4)`,
			e.ID, ref.Type, ref.Reference, ref.FullAmount)
		if err != nil {
			return errors.Wrap(err, "expenses.saveExternalRecords")
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
		return errors.Wrap(err, "expenses.addExternalRecords")
	}
	for rows.Next() {
		var typeValue, reference string
		var oldamount int64
		err = rows.Scan(&typeValue, &reference, &oldamount)
		if err != nil {
			return errors.Wrap(err, "expenses.addExternalRecords")
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
	return errors.Wrap(err, "expenses.deleteExpense")
}

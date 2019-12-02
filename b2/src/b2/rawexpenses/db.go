package rawexpenses 

import "database/sql"
import "errors"

func findRawExpenses(db *sql.DB) ([]uint64, error) {
    rows, err := db.Query("select rid from rawdata")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var rids []uint64
    for rows.Next() {
        var eid uint64
        err = rows.Scan(&eid)
        if err != nil {
            return nil, err
        }
        rids = append(rids, eid)
    }
    return rids, err
}

func loadRawExpense(rid uint64, db *sql.DB) (*RawExpense, error) {
    rows, err := db.Query(`
        select
            r.rawStr,
            r.importDate,
            r.aid
        from
            rawdata r
        where
            r.rid = $1`,
            rid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    rawexpense := new(RawExpense)
    if rows.Next() {
        err = rows.Scan(&rawexpense.Data,
                        &rawexpense.Date,
                        &rawexpense.AccountID)
        rawexpense.ID = rid
    } else {
        return nil, errors.New("404")
    }
    if err != nil {
        return nil, err
    }
    return rawexpense, nil
}

func createRawExpense(e *RawExpense, db *sql.DB) error {
    // todo: check values are legit before writing
    res, err := db.Exec(`insert into
                            rawdata
                                (rawStr,
                                 importDate,
                                 aid)
                             values
                                 ($1,
                                  $2,
                                  $3);`,
                          e.Data,
                          e.Date,
                          e.AccountID)
    if err != nil {
        return err
    }
    rid, err := res.LastInsertId()
    if (err == nil || rid < 1) {
        e.ID = uint64(rid)
    } else {
        return errors.New("Error creating new rawexpense")
    }

    return err
}

func updateRawExpense(e *RawExpense, db *sql.DB) error {
    e.RLock()
    defer e.RUnlock()
    // Todo: Check values are legit before writing
    //_, err := db.Exec("update expenses set aid = $1, description = $2, amount = $3, ccy = $4, amountFX = $5, ccyFX = $6, fxRate = $7, commission = $8, date = $9, temporary = $10, reference = $11, detaileddescription = $12, processDate = $13 where eid = $14; delete from classifications where eid = $15; insert into classifications  (eid, cid, confirmed) values ($16, $17, $18)", e.AccountID, e.Description, e.Amount, e.Currency, e.FX.Amount, e.FX.Currency, e.FX.Rate, e.Commission, e.Date, e.Metadata.Temporary, e.TransactionReference, e.DetailedDescription, e.ProcessDate, e.ID, e.ID, e.ID, e.Metadata.Classification, e.Metadata.Confirmed)
    //if err != nil {
      //  return err
    //}
    //return err
    return nil
}

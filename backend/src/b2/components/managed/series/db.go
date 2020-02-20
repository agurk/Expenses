package series

import (
	"b2/errors"
	"database/sql"
	"fmt"
)

func loadSeries(sid uint64, db *sql.DB) (*Series, error) {
	rows, err := db.Query(`
        select
			sid,
			asid,
			date,
			amountwhole,
			amountfractional	
        from
			assetseries
        where
            sid = $1`,
		sid)
	if err != nil {
		return nil, errors.Wrap(err, "series.loadSeries")
	}
	defer rows.Close()
	series := new(Series)
	if rows.Next() {
		err = rows.Scan(&series.ID,
			&series.AssetID,
			&series.Date,
			&series.WholeAmount,
			&series.FractionalAmount)
	} else {
		return nil, errors.New("Series not found", errors.ThingNotFound, "series.loadSeries", true)
	}
	if err != nil {
		return nil, errors.Wrap(err, "series.loadSeries")
	}
	return series, nil
}

func findSeries(query *Query, db *sql.DB) ([]uint64, error) {
	var args []interface{}
	var where bool
	dbQuery := `
		select
			sid
		from
			assetseries `
	if query.AssetID > 0 {
		dbQuery += `
		where
			asid = $1 `
		args = append(args, query.AssetID)
		where = true
	}
	if query.Date != "" {
		if !where {
			dbQuery += " where "
			where = true
		} else {
			dbQuery += " and "
		}
		args = append(args, query.Date)
		dbQuery += fmt.Sprintf(" date <= $%d ", len(args))
	}
	if query.OnlyLatest {
		dbQuery += ` order by date desc limit 1 `
	}
	rows, err := db.Query(dbQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "series.findSeries")
	}
	defer rows.Close()
	var sids []uint64
	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid)
		if err != nil {
			return nil, errors.Wrap(err, "series.findSeries")
		}
		sids = append(sids, aid)
	}
	return sids, errors.Wrap(err, "series.findSeries")
}

func createSeries(series *Series, db *sql.DB) error {
	series.Lock()
	defer series.Unlock()
	res, err := db.Exec(`insert into
							assetseries (
								asid,
								date,
								amountwhole,
								amountfractional)
							values ($1, $2, $3, $4)`,
		series.AssetID,
		series.Date,
		series.WholeAmount,
		series.FractionalAmount)

	if err != nil {
		return errors.Wrap(err, "series.createSeries")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		series.ID = uint64(rid)
	} else {
		return errors.New("Error creating new series", errors.InternalError, "series.createSeries", false)
	}
	return errors.Wrap(err, "series.createSeries")
}

func updateSeries(series *Series, db *sql.DB) error {
	series.RLock()
	defer series.RUnlock()
	_, err := db.Exec(`
		update
			assetseries
		set
			asid = $1,
			date = $2,
			amountwhole = $3,
			amountfractional = $4
		where
			sid = $5`,
		series.AssetID,
		series.Date,
		series.WholeAmount,
		series.FractionalAmount,
		series.ID)
	return errors.Wrap(err, "series.updateSeries")
}

func deleteSeries(series *Series, db *sql.DB) error {
	_, err := db.Exec(`
        delete from
			assetseries
        where
            sid = $1`,
		series.ID)
	return errors.Wrap(err, "series.deleteSeries(delete)")

}

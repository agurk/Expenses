package assets

import (
	"b2/errors"
	"database/sql"
)

func loadAsset(asid uint64, db *sql.DB) (*Asset, error) {
	rows, err := db.Query(`
        select
            asid,
            name,
			reference,
			type,
			symbol
        from
			assets
        where
            asid = $1`,
		asid)
	if err != nil {
		return nil, errors.Wrap(err, "assets.loadAsset")
	}
	defer rows.Close()
	asset := new(Asset)
	if rows.Next() {
		err = rows.Scan(&asset.ID,
			&asset.Name,
			&asset.Reference,
			&asset.Variety,
			&asset.Symbol)
	} else {
		return nil, errors.New("Asset not found", errors.ThingNotFound, "assets.loadAsset", true)
	}
	if err != nil {
		return nil, errors.Wrap(err, "assets.loadAsset")
	}
	return asset, nil
}

func findAssets(db *sql.DB) ([]uint64, error) {
	rows, err := db.Query("select asid from assets")
	if err != nil {
		return nil, errors.Wrap(err, "assets.findAssets")
	}
	defer rows.Close()
	var aids []uint64
	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid)
		if err != nil {
			return nil, errors.Wrap(err, "assets.findAssets")
		}
		aids = append(aids, aid)
	}
	return aids, errors.Wrap(err, "assets.findAssets")
}

func createAsset(asset *Asset, db *sql.DB) error {
	asset.Lock()
	defer asset.Unlock()
	res, err := db.Exec(`insert into
							assets (
								name,
								reference,
								type,
								symbol)
							values ($1, $2, $3, $4)`,
		asset.Name,
		asset.Reference,
		asset.Variety,
		asset.Symbol)

	if err != nil {
		return errors.Wrap(err, "assets.createAsset")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		asset.ID = uint64(rid)
	} else {
		return errors.New("Error creating new asset", errors.InternalError, "assets.createAsset", false)
	}
	return errors.Wrap(err, "assets.createAsset")
}

func updateAsset(asset *Asset, db *sql.DB) error {
	asset.RLock()
	defer asset.RUnlock()
	_, err := db.Exec(`
		update
			assets
		set
			name = $1,
			reference = $2,
			type = $3,
			symbol = $4
		where
			aid = $5`,
		asset.Name,
		asset.Reference,
		asset.Variety,
		asset.Symbol,
		asset.ID)
	return errors.Wrap(err, "assets.updateAsset")
}

func deleteAsset(asset *Asset, db *sql.DB) error {
	_, err := db.Exec(`
        delete from
			assets
        where
            asid = $1`,
		asset.ID)
	return errors.Wrap(err, "assets.deleteAsset(delete)")

}

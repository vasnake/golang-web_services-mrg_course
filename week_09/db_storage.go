package main

import (
	"database/sql"
	"fmt"
)

type Photo struct {
	ID     int
	UserID int
	Path   string
}

// storage implementation

type StDb struct {
	db *sql.DB
}

// storage factory
func NewDbStorage(db *sql.DB) *StDb {
	return &StDb{
		db: db,
	}
}

// insert (user, path) to table
func (st *StDb) Add(p *Photo) error {
	res, err := st.db.Exec("INSERT INTO photos(user_id, path) VALUES(?, ?)",
		p.UserID, p.Path)
	if err != nil {
		return err
	}

	li, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if li == 0 {
		return fmt.Errorf("no last insert id")
	}

	return nil
}

// list user files
func (st *StDb) GetPhotos(userID int) ([]*Photo, error) {
	photos := make([]*Photo, 0, 10)

	rows, err := st.db.Query("select id, user_id, path from photos where user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Photo{}
		err := rows.Scan(&item.ID, &item.UserID, &item.Path)
		if err != nil {
			return nil, err
		}

		photos = append(photos, item)
	}

	return photos, nil
}

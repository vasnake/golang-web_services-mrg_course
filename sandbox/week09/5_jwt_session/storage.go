package jwt_session

import (
	"database/sql"
	"fmt"
)

// same shit

type Photo struct {
	ID      int
	UserID  uint32
	Path    string
	Comment string
	Rating  int64
}

// Storage implementation (add, list, rate)
type StDb struct {
	db *sql.DB
}

func NewDbStorage(db *sql.DB) *StDb {
	return &StDb{
		db: db,
	}
}

func (st *StDb) Add(p *Photo) error {
	res, err := st.db.Exec("INSERT INTO photos(user_id, path, comment) VALUES(?, ?, ?)",
		p.UserID, p.Path, p.Comment)
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

func (st *StDb) GetPhotos(userID uint32) ([]*Photo, error) {
	photos := make([]*Photo, 0, 10)

	rows, err := st.db.Query("select id, user_id, path, comment, rating from photos where user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := &Photo{}
		err := rows.Scan(&item.ID, &item.UserID, &item.Path, &item.Comment, &item.Rating)
		if err != nil {
			return nil, err
		}
		photos = append(photos, item)
	}

	return photos, nil
}

func (st *StDb) Rate(photoID uint32, rate int) error {
	_, err := st.db.Exec("UPDATE photos SET rating = rating + ? WHERE id = ?", rate, photoID)
	return err
}

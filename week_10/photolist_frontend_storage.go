package main

import (
	"database/sql"
	"fmt"
)

type Photo struct {
	ID     int
	UserID uint32
	Path   string
	// new fields
	Comment string
	Rating  int64
	Liked   bool // null to bool conversion
}

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

// photos left join likes, with null as liked if user don't click like
func (st *StDb) GetPhotos(userID uint32) ([]*Photo, error) {
	photos := make([]*Photo, 0, 10)

	rows, err := st.db.Query(`SELECT id, photos.user_id, path, comment, rating, user_photos_likes.photo_id AS liked
	FROM photos
	LEFT JOIN user_photos_likes ON photos.id=user_photos_likes.photo_id
	WHERE photos.user_id = ?
	ORDER BY id DESC`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Photo{}
		var liked sql.NullInt64
		err := rows.Scan(&item.ID, &item.UserID, &item.Path, &item.Comment, &item.Rating, &liked)
		if err != nil {
			return nil, err
		}

		item.Liked = liked.Valid // null to bool conversion

		photos = append(photos, item)
	}

	return photos, nil
}

// insert or delete record in likes; update rating in photos
func (st *StDb) Rate(photoID uint32, userID uint32, rate int) error {
	var res sql.Result
	var err error

	if rate >= 0 {
		res, err = st.db.Exec(`INSERT IGNORE INTO user_photos_likes(photo_id, user_id) VALUES(?, ?)`, photoID, userID)
	} else {
		res, err = st.db.Exec(`DELETE FROM user_photos_likes WHERE photo_id = ? AND user_id = ?`, photoID, userID)
	}

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	// dont update rating twice, if inserted or deleted already
	if aff <= 0 {
		return nil
	}

	_, err = st.db.Exec("UPDATE photos SET rating = rating + ? WHERE id = ?", rate, photoID)
	return err
}

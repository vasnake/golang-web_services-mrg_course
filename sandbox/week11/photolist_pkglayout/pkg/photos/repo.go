package photos

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Photo struct {
	ID     uint32 `json:"id"`
	UserID uint32 `json:"-"`
	// User     *User  `json:"user"`
	URL     string `json:"url"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	Liked   bool   `json:"liked"`
}

func (ph *Photo) Id() string {
	return strconv.Itoa(int(ph.ID))
}

type StDb struct {
	db *sql.DB
}

func NewDbStorage(db *sql.DB) *StDb {
	return &StDb{
		db: db,
	}
}

func (st *StDb) Add(p *Photo) (uint32, error) {
	res, err := st.db.Exec("INSERT INTO photos(user_id, path, comment) VALUES(?, ?, ?)",
		p.UserID, p.URL, p.Comment)
	if err != nil {
		return 0, err
	}

	li, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if li == 0 {
		return 0, fmt.Errorf("no last insert id")
	}

	return uint32(li), nil
}

func (st *StDb) GetByID(photoID, currentUserID uint32) (*Photo, error) {
	rows := st.db.QueryRow(`SELECT
		photos.id as id, photos.user_id, path, comment, rating, 
		user_photos_likes.photo_id as is_liked
	   FROM photos 
	   LEFT JOIN users ON photos.user_id=users.id
	   LEFT JOIN user_photos_likes ON user_photos_likes.photo_id=photos.id and user_photos_likes.user_id = ?
	   WHERE photos.id = ?
	   ORDER BY id DESC`, currentUserID, photoID)
	item := &Photo{}
	var isLiked sql.NullInt64
	err := rows.Scan(&item.ID, &item.UserID, &item.URL, &item.Comment, &item.Rating, &isLiked)
	if err != nil {
		return nil, err
	}
	item.Liked = isLiked.Valid
	return item, nil
}

func (st *StDb) GetPhotos(userID, currentUserID uint32) ([]*Photo, error) {
	photos := make([]*Photo, 0, 10)

	// TODO add limit, offset
	rows, err := st.db.Query(`SELECT 
	photos.id as id, photos.user_id, path, comment, rating, 
		   users.login as user_login, 
		   user_photos_likes.photo_id as is_liked, 
		   user_follows.follow_id as is_followed
	   FROM photos 
	   LEFT JOIN users ON photos.user_id=users.id
	   LEFT JOIN user_photos_likes ON user_photos_likes.photo_id=photos.id and user_photos_likes.user_id = ?
	   LEFT JOIN user_follows ON  user_follows.follow_id=photos.user_id and user_follows.user_id = ?
	   WHERE photos.user_id = ?
	   ORDER BY id DESC`, currentUserID, currentUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Photo{}
		var isLiked, isFollowed sql.NullInt64
		var userLogin string
		err := rows.Scan(&item.ID, &item.UserID, &item.URL, &item.Comment, &item.Rating, &userLogin, &isLiked, &isFollowed)
		if err != nil {
			return nil, err
		}
		item.Liked = isLiked.Valid
		// item.Followed = isFollowed.Valid
		photos = append(photos, item)
	}

	return photos, nil
}

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
	// dont update rating twice
	if aff <= 0 {
		return nil
	}
	_, err = st.db.Exec("UPDATE photos SET rating = rating + ? WHERE id = ?", rate, photoID)
	return err
}

package user

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	errUserNotFound = errors.New("No user record found")
	errBadPass      = errors.New("Bad password")
	errUserExists   = errors.New("User Exists")
)

type UserRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) LookupByIDs(currUserID uint32, ids []uint32) ([]*User, []error) {
	// fortunately for me - almost direct copy-paste from https://gqlgen.com/reference/dataloaders/

	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids)+1)
	args = append(args, currUserID)
	for i := 0; i < len(ids); i++ {
		placeholders[i] = "?"
		args = append(args, ids[i])
	}

	q := `SELECT id, login, user_follows.follow_id FROM users 
	LEFT JOIN user_follows ON user_follows.follow_id=users.id and user_follows.user_id = ?
	WHERE users.id IN (` + strings.Join(placeholders, ",") + ")"
	res, err := repo.db.Query(q, args...)
	if err != nil {
		log.Println("LookupByIDs query err:", err)
		return nil, []error{err}
	}
	defer res.Close()

	users := make(map[uint32]*User, len(ids))
	for res.Next() {
		user := &User{}
		var isFollowed sql.NullInt64
		err := res.Scan(&user.ID, &user.Login, &isFollowed)
		if err != nil {
			return nil, []error{err}
		}
		user.Followed = &isFollowed.Valid
		users[user.ID] = user
	}

	output := make([]*User, len(ids))
	for i, id := range ids {
		output[i] = users[id]
	}
	return output, nil
}

func (repo *UserRepository) Create(login, email, passIn string) (*User, error) {
	salt := makeSalt(8)
	pass := repo.hashPass(passIn, salt)

	user := &User{
		ID:    0,
		Ver:   0,
		Email: email,
	}

	err := repo.db.QueryRow("SELECT id, ver, login FROM users WHERE email = ? OR login = ?", email, login).
		Scan(&user.ID, &user.Ver, &user.Login)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("db error: %v", err)
	}
	if err != sql.ErrNoRows {
		return user, errUserExists
	}

	result, err := repo.db.Exec("INSERT INTO users(login, email, password) VALUES(?, ?, ?)", login, email, pass)
	if err != nil {
		return nil, fmt.Errorf("insert error: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("no rows affected")
	}
	uid, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("LastInsertId err: %v", err)
	}
	user.ID = uint32(uid)

	return user, nil
}

func (repo *UserRepository) passwordIsValid(pass string, row *sql.Row) (*User, error) {
	var (
		dbPass []byte
		user   = &User{}
	)
	err := row.Scan(&user.ID, &user.Login, &user.Ver, &dbPass)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, err
	}

	salt := dbPass[0:8]
	if !bytes.Equal(repo.hashPass(pass, salt), dbPass) {
		return nil, errBadPass
	}
	return user, nil
}

func (repo *UserRepository) GetByLogin(login string) (*User, error) {
	row := repo.db.QueryRow("SELECT id, login, email, ver FROM users WHERE login = ?", login)
	return parseRowToUser(row)
}

func (repo *UserRepository) GetByID(id uint32) (*User, error) {
	row := repo.db.QueryRow("SELECT id, login, email, ver FROM users WHERE id = ?", id)
	return parseRowToUser(row)
}

func (repo *UserRepository) CheckPasswordByUserID(uid uint32, pass string) (*User, error) {
	row := repo.db.QueryRow("SELECT id, login, ver, password FROM users WHERE id = ?", uid)
	return repo.passwordIsValid(pass, row)
}

func (repo *UserRepository) CheckPasswordByLogin(login, pass string) (*User, error) {
	row := repo.db.QueryRow("SELECT id, login, ver, password FROM users WHERE login = ?", login)
	return repo.passwordIsValid(pass, row)
}

func (repo *UserRepository) UpdatePassword(userID uint32, pass string) error {
	salt := makeSalt(8)
	passHash := repo.hashPass(pass, salt)
	_, err := repo.db.Exec("UPDATE users SET password = ?, ver = ver + 1 WHERE id = ?",
		passHash, userID)
	return err
}

func (repo *UserRepository) hashPass(plainPassword string, salt []byte) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}

func makeSalt(n int) []byte {
	salt := make([]byte, 8)
	rand.Read(salt)
	return salt
}

func parseRowToUser(row *sql.Row) (*User, error) {
	user := &User{}
	err := row.Scan(&user.ID, &user.Login, &user.Email, &user.Ver)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Follow(userID uint32, currentUserID uint32, rate int) error {
	var res sql.Result
	var err error
	if rate == 1 {
		res, err = repo.db.Exec(`INSERT IGNORE INTO user_follows(follow_id, user_id) VALUES(?, ?)`, userID, currentUserID)
	} else {
		res, err = repo.db.Exec(`DELETE FROM user_follows WHERE follow_id = ? AND user_id = ?`, userID, currentUserID)
	}
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	// dont update rating twice
	if aff <= 0 {
		return nil
	}
	_, err = repo.db.Exec("UPDATE users SET followers_cnt = followers_cnt + ? WHERE id = ?", rate, userID)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec("UPDATE users SET following_cnt = following_cnt + ? WHERE id = ?", rate, currentUserID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) IsFollowed(userID, currUserID uint32) (bool, error) {
	log.Println("call UserRepository.IsFollowed - maybe user dataloader?", userID, currUserID)
	q := `SELECT count(*) as cnt FROM user_follows 
		WHERE user_id = ? AND follow_id = ?`
	var cnt uint32
	err := repo.db.QueryRow(q, currUserID, userID).Scan(&cnt)
	return cnt != 0, err
}

func (repo *UserRepository) GetFollowedUsers(userID uint32) ([]*User, error) {
	// TODO add limit, offset
	rows, err := repo.db.Query(`SELECT users.id, users.login 
	FROM user_follows 
	LEFT JOIN users ON users.id = user_follows.follow_id
	WHERE user_follows.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*User, 0, 10)
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Login)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (repo *UserRepository) GetRecomendedUsers(userID uint32) ([]*User, error) {
	// TODO add limit, offset
	rows, err := repo.db.Query(`select users.id, users.login 
	from users 
	left join user_follows on users.id = user_follows.follow_id and user_follows.user_id = ?
	where users.id != ? and user_follows.user_id is null`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*User, 0, 10)
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Login)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func IsErrUserNotFound(err error) bool {
	return err == errUserNotFound
}

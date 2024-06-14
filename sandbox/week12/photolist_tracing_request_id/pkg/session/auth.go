package session

import (
	"context"
	"database/sql"

	"photolist/pkg/utils/randutils"
)

type AuthService struct {
	DB *sql.DB
}

func (as *AuthService) Check(ctx context.Context, c *AuthCheckIn) (*AuthSession, error) {
	sess := &AuthSession{}
	row := as.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, c.GetSessKey())
	err := row.Scan(&sess.UserID)
	if err == sql.ErrNoRows {
		// log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		// log.Println("CheckSession err:", err)
		return nil, err
	}
	return sess, nil
}

func (as *AuthService) Create(ctx context.Context, u *AuthUserIn) (*AuthSession, error) {
	sessID := randutils.RandStringRunes(32)
	_, err := as.DB.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, u.GetUserID())
	if err != nil {
		return nil, err
	}
	return &AuthSession{
		ID:     sessID,
		UserID: u.GetUserID(),
	}, nil
}

func (as *AuthService) DestroyCurrent(ctx context.Context, s *AuthSession) (*AuthNothing, error) {
	_, err := as.DB.Exec("DELETE FROM sessions WHERE id = ?",
		s.GetID())
	return &AuthNothing{}, err
}

func (as *AuthService) DestroyAll(ctx context.Context, u *AuthUserIn) (*AuthNothing, error) {
	_, err := as.DB.Exec("DELETE FROM sessions WHERE user_id = ?",
		u.GetUserID())
	return &AuthNothing{}, err
}

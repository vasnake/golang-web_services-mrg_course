package sql_storage

import (
	"fmt"
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/*
	go test -v -run Storage -coverprofile=db_st.out && go tool cover -html=db_st.out -o db_st.html && rm db_st.out
*/

func TestStorageAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	st := NewDbStorage(db)

	userID := 1
	path := "test"
	testPhoto := &Photo{UserID: userID, Path: path}

	//ok query
	mock.
		ExpectExec(`INSERT INTO photos`).
		WithArgs(userID, path).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = st.Add(testPhoto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO photos`).
		WithArgs(userID, path).
		WillReturnError(fmt.Errorf("bad query"))

	err = st.Add(testPhoto)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO photos`).
		WithArgs(userID, path).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	err = st.Add(testPhoto)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// last id error
	mock.
		ExpectExec(`INSERT INTO photos`).
		WithArgs(userID, path).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = st.Add(testPhoto)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestStorageGetPhotos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "user_id", "path"})
	expect := []*Photo{
		{1, userID, "tree"},
		{2, userID, "minion"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.UserID, item.Path)
	}

	mock.
		ExpectQuery("select id, user_id, path from photos where").
		WithArgs(userID).
		WillReturnRows(rows)

	st := NewDbStorage(db)
	items, err := st.GetPhotos(userID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(expect, items) {
		t.Errorf("results not match")
		return
	}

	// query error
	mock.
		ExpectQuery("select id, user_id, path from photos where").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = st.GetPhotos(userID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "path"}).
		AddRow(1, "camera")

	mock.
		ExpectQuery("select id, user_id, path from photos where").
		WithArgs(userID).
		WillReturnRows(rows)

	_, err = st.GetPhotos(userID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

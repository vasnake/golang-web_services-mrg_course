package main

import (
	"bytes"
	"fmt"
	"log"

	// "fmt"
	"html/template"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

// mock storage
/*
	находясь в папке с этим файлом handlers_test.go
	mockgen -source=storage.go -destination=handlers_mock.go -package=main Storage

	go test -v -run Handler -coverprofile=handler.out && go tool cover -html=handler.out -o handler.html && rm handler.out
*/

func TestHandlerGetPhotos(t *testing.T) {

	log.SetOutput(ioutil.Discard) // disable log

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	// storage mock
	st := NewMockStorage(ctrl)

	// create service to test
	service := &PhotolistHandler{
		St:   st,
		Tmpl: NewTemplates(),
	}

	// mock result
	resultItems := []*Photo{
		{1, 1, "my_photo_name"},
	}

	// actual mock
	// тут мы записываем последовательность вызовов и результат
	st.EXPECT().GetPhotos(0).Return(resultItems, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// call tested method, happy path
	service.List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	img := `"/images/my_photo_name_160.jpg"`
	if !bytes.Contains(body, []byte(img)) {
		t.Errorf("no image found")
		return
	}

	// check GetPhotos error case

	// тут мы записываем последовательность вызовов и результат
	st.EXPECT().GetPhotos(0).Return(nil, fmt.Errorf("no results"))

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	// call with mock
	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	// template expand error case

	// template mock
	service.Tmpl, _ = template.New("tmplError").Parse("{{.NotExist}}")

	// storage mock
	st.EXPECT().GetPhotos(0).Return(resultItems, nil)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	// call with mocks
	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

}

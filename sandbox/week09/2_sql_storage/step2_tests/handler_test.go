package sql_storage

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

/*
	находять в папке с этим файлом handlers_test.go
	mockgen -source=handlers.go -destination=handlers_mock.go -package=main Storage
	go test -v -run Handler -coverprofile=handler.out && go tool cover -html=handler.out -o handler.html && rm handler.out
*/

func TestHandlerGetPhotos(t *testing.T) {

	log.SetOutput(ioutil.Discard)

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockStorage(ctrl)
	service := &PhotolistHandler{
		St:   st,
		Tmpl: NewTemplates(),
	}

	resultItems := []*Photo{
		{1, 1, "my_photo_name"},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetPhotos(0).Return(resultItems, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	service.List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	img := `"/images/my_photo_name_160.jpg"`
	if !bytes.Contains(body, []byte(img)) {
		t.Errorf("no image found")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetPhotos(0).Return(nil, fmt.Errorf("no results"))

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	// template expand error
	service.Tmpl, _ = template.New("tmplError").Parse("{{.NotExist}}")

	st.EXPECT().GetPhotos(0).Return(resultItems, nil)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

}

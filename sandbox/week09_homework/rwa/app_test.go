package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mcuadros/go-lookup"
	"gopkg.in/d4l3k/messagediff.v1"
)

type FakeTime struct {
	Valid bool
}

func (ft *FakeTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	_, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	// fmt.Println(string(data), tt)
	ft.Valid = err == nil
	return err
}

type TestProfile struct {
	ID        string   `json:"id" testdiff:"ignore"`
	Email     string   `json:"email"`
	CreatedAt FakeTime `json:"createdAt"`
	UpdatedAt FakeTime `json:"updatedAt"`
	Username  string   `json:"username"`
	Bio       string   `json:"bio"`
	Image     string   `json:"image"`
	Token     string   `json:"token" testdiff:"ignore"`
	Following bool
}

type TestArticle struct {
	Author         TestProfile `json:"author"`
	Body           string      `json:"body"`
	CreatedAt      FakeTime    `json:"createdAt"`
	Description    string      `json:"description"`
	Favorited      bool        `json:"favorited"`
	FavoritesCount int         `json:"favoritesCount"`
	Slug           string      `json:"slug" testdiff:"ignore"`
	TagList        []string    `json:"tagList"`
	Title          string      `json:"title"`
	UpdatedAt      FakeTime    `json:"updatedAt"`
}

func strP(in string) *string {
	return &in
}

type ApiTestCase struct {
	Name           string
	Event          []string
	Method         string
	Body           string
	URL            string
	TokenName      string
	ResponseStatus int
	ResponsePath   string
	Expected       interface{}
	Before         func()
	After          func(*http.Response, []byte, interface{}) error
}

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

func WeirdMagicClone(in interface{}) interface{} {
	return reflect.New(reflect.TypeOf(in).Elem()).Interface()
}

func TestApp(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var (
		app = GetApp()
		ts  = httptest.NewServer(app)

		// username = RandStringRunes(16)
		username = "golang"
		// apiurl   = "https://conduit.productionready.io/api"
		apiurl   = ts.URL + "/api"
		password = "love"
	)

	tplParams := map[string]string{
		"APIURL":   apiurl,
		"EMAIL":    username + "@example.com",
		"PASSWORD": password,
		"USERNAME": username,
		"BIO":      "Info about " + username,

		"EMAIL2":    username + "_second@example.com",
		"USERNAME2": username + "_second",
	}

	replaceRe := regexp.MustCompile("{{(.*?)}}")
	replaceBrackets := strings.NewReplacer("{", "", "}", "")
	replacer := func(key []byte) []byte {
		k := replaceBrackets.Replace(string(key))
		val, ok := tplParams[k]
		if !ok {
			t.Fatalf("not found key %s during tpl substitution", string(key))
		}
		return []byte(val)
	}

	testCases := []*ApiTestCase{
		&ApiTestCase{
			Name:           "Auth - Register",
			Method:         "POST",
			Body:           "{\"user\":{\"email\":\"{{EMAIL}}\", \"password\":\"{{PASSWORD}}\", \"username\":\"{{USERNAME}}\"}}",
			URL:            "{{APIURL}}/users",
			ResponseStatus: 201,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME"],
					},
				}
			},
			Before: nil,
			After:  nil,
		},
		&ApiTestCase{
			Name:           "Auth - Login",
			Method:         "POST",
			Body:           "{\"user\":{\"email\":\"{{EMAIL}}\", \"password\":\"{{PASSWORD}}\"}}",
			URL:            "{{APIURL}}/users/login",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME"],
					},
				}
			},
			Before: nil,
			After: func(r *http.Response, body []byte, resp interface{}) error {
				val, err := lookup.LookupString(resp, "User.Token")
				if err != nil {
					return err
				}
				tplParams["token1"] = val.String()
				return nil
			},
		},
		&ApiTestCase{
			Name:           "Auth - Current User",
			Method:         "GET",
			Body:           "",
			URL:            "{{APIURL}}/user",
			TokenName:      "token1",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME"],
					},
				}
			},
			Before: nil,
			After:  nil,
		},
		&ApiTestCase{
			Name:           "Auth - Update User",
			Method:         "PUT",
			Body:           `{"user":{"email":"{{EMAIL}}","bio":"{{BIO}}"}}`,
			URL:            "{{APIURL}}/user",
			TokenName:      "token1",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME"],
						Bio:       tplParams["BIO"],
					},
				}
			},
			Before: func() {
				tplParams["EMAIL"] = "u_" + tplParams["EMAIL"]
			},
			After: func(r *http.Response, body []byte, resp interface{}) error {
				val, err := lookup.LookupString(resp, "User.Token")
				if err != nil {
					return err
				}
				tplParams["token1"] = val.String()
				return nil
			},
		},
		&ApiTestCase{
			Name:           "Auth - Current User after Update",
			Method:         "GET",
			URL:            "{{APIURL}}/user",
			TokenName:      "token1",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME"],
						Bio:       tplParams["BIO"],
					},
				}
			},
		},
		&ApiTestCase{
			Name:           "Auth - Register second user",
			Method:         "POST",
			Body:           "{\"user\":{\"email\":\"{{EMAIL2}}\", \"password\":\"{{PASSWORD}}\", \"username\":\"{{USERNAME2}}\"}}",
			URL:            "{{APIURL}}/users",
			ResponseStatus: 201,
			Expected: func() interface{} {
				return &struct {
					User TestProfile
				}{
					User: TestProfile{
						Email:     tplParams["EMAIL2"],
						CreatedAt: FakeTime{true},
						UpdatedAt: FakeTime{true},
						Username:  tplParams["USERNAME2"],
					},
				}
			},
			After: func(r *http.Response, body []byte, resp interface{}) error {
				val, err := lookup.LookupString(resp, "User.Token")
				if err != nil {
					return err
				}
				tplParams["token2"] = val.String()
				return nil
			},
		},

		&ApiTestCase{
			Name:           "Articles - Create Article - First user",
			Method:         "POST",
			Body:           `{"article":{"title":"How to write golang tests", "description":"I have problem with mondodb mocking", "body":"Any ideas how to write some intermidiate layer atop collection?", "tagList":["golang","testing", "gomock"]}}`,
			URL:            "{{APIURL}}/articles",
			TokenName:      "token1",
			ResponseStatus: 201,
			Expected: func() interface{} {
				return &struct {
					Article TestArticle
				}{
					Article: TestArticle{
						Author: TestProfile{
							Bio:      tplParams["BIO"],
							Username: tplParams["USERNAME"],
						},
						Body:        "Any ideas how to write some intermidiate layer atop collection?",
						Title:       "How to write golang tests",
						Description: "I have problem with mondodb mocking",
						CreatedAt:   FakeTime{true},
						UpdatedAt:   FakeTime{true},
						TagList:     []string{"golang", "testing", "gomock"},
					},
				}
			},
			After: func(r *http.Response, body []byte, resp interface{}) error {
				val, err := lookup.LookupString(resp, "Article.Slug")
				if err != nil {
					return err
				}
				tplParams["slug1"] = val.String()
				return nil
			},
		},
		&ApiTestCase{
			Name:           "Articles - Create Article - Second user",
			Method:         "POST",
			Body:           `{"article":{"title":"What will be released first, Half-Life 3 or 3-rd part of golang course?", "description":"Who knows topics in new course?", "body":"Will we use JWT-tokens in homework?", "tagList":["halflife3","coursera"]}}`,
			URL:            "{{APIURL}}/articles",
			TokenName:      "token2",
			ResponseStatus: 201,
			Expected: func() interface{} {
				return &struct {
					Article TestArticle
				}{
					Article: TestArticle{
						Author: TestProfile{
							Username: tplParams["USERNAME2"],
						},
						Body:        "Will we use JWT-tokens in homework?",
						Title:       "What will be released first, Half-Life 3 or 3-rd part of golang course?",
						Description: "Who knows topics in new course?",
						CreatedAt:   FakeTime{true},
						UpdatedAt:   FakeTime{true},
						TagList:     []string{"halflife3", "coursera"},
					},
				}
			},
			After: func(r *http.Response, body []byte, resp interface{}) error {
				val, err := lookup.LookupString(resp, "Article.Slug")
				if err != nil {
					return err
				}
				tplParams["slug2"] = val.String()
				return nil
			},
		},

		&ApiTestCase{
			Name:           "Articles - All Articles",
			Method:         "GET",
			URL:            "{{APIURL}}/articles",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					Articles      []TestArticle `json:"articles"`
					ArticlesCount int           `json:"articlesCount"`
				}{
					Articles: []TestArticle{
						TestArticle{
							Slug: tplParams["slug1"],
							Author: TestProfile{
								Bio:      tplParams["BIO"],
								Username: tplParams["USERNAME"],
							},
							Body:        "Any ideas how to write some intermidiate layer atop collection?",
							Title:       "How to write golang tests",
							Description: "I have problem with mondodb mocking",
							CreatedAt:   FakeTime{true},
							UpdatedAt:   FakeTime{true},
							TagList:     []string{"golang", "testing", "gomock"},
						},
						TestArticle{
							Slug: tplParams["slug2"],
							Author: TestProfile{
								Username: tplParams["USERNAME2"],
							},
							Body:        "Will we use JWT-tokens in homework?",
							Title:       "What will be released first, Half-Life 3 or 3-rd part of golang course?",
							Description: "Who knows topics in new course?",
							CreatedAt:   FakeTime{true},
							UpdatedAt:   FakeTime{true},
							TagList:     []string{"halflife3", "coursera"},
						},
					},
					ArticlesCount: 2,
				}
			},
			Before: nil,
			After:  nil,
		},

		&ApiTestCase{
			Name:           "Articles - by author",
			Method:         "GET",
			URL:            "{{APIURL}}/articles?author={{USERNAME2}}",
			TokenName:      "token1",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					Articles      []TestArticle `json:"articles"`
					ArticlesCount int           `json:"articlesCount"`
				}{
					Articles: []TestArticle{
						TestArticle{
							Slug: tplParams["slug2"],
							Author: TestProfile{
								Username: tplParams["USERNAME2"],
							},
							Body:        "Will we use JWT-tokens in homework?",
							Title:       "What will be released first, Half-Life 3 or 3-rd part of golang course?",
							Description: "Who knows topics in new course?",
							CreatedAt:   FakeTime{true},
							UpdatedAt:   FakeTime{true},
							TagList:     []string{"halflife3", "coursera"},
						},
					},
					ArticlesCount: 1,
				}
			},
			Before: nil,
			After:  nil,
		},
		&ApiTestCase{
			Name:           "Articles - by tag",
			Method:         "GET",
			URL:            "{{APIURL}}/articles?tag=halflife3",
			TokenName:      "token1",
			ResponseStatus: 200,
			Expected: func() interface{} {
				return &struct {
					Articles      []TestArticle `json:"articles"`
					ArticlesCount int           `json:"articlesCount"`
				}{
					Articles: []TestArticle{
						TestArticle{
							Slug: tplParams["slug2"],
							Author: TestProfile{
								Username: tplParams["USERNAME2"],
							},
							Body:        "Will we use JWT-tokens in homework?",
							Title:       "What will be released first, Half-Life 3 or 3-rd part of golang course?",
							Description: "Who knows topics in new course?",
							CreatedAt:   FakeTime{true},
							UpdatedAt:   FakeTime{true},
							TagList:     []string{"halflife3", "coursera"},
						},
					},
					ArticlesCount: 1,
				}
			},
			Before: nil,
			After:  nil,
		},

		&ApiTestCase{
			Name:           "No Auth - Current User - No Auth",
			Method:         "GET",
			URL:            "{{APIURL}}/user",
			TokenName:      "", // none
			ResponseStatus: 401,
		},
		&ApiTestCase{
			Name:           "No Auth - Current User Logout - Require Auth",
			Method:         "POST",
			URL:            "{{APIURL}}/user/logout",
			TokenName:      "", // none
			ResponseStatus: 401,
		},
		&ApiTestCase{
			Name:           "No Auth - Current User Logout",
			Method:         "POST",
			URL:            "{{APIURL}}/user/logout",
			TokenName:      "token1",
			ResponseStatus: 200,
		},
		&ApiTestCase{
			Name:           "No Auth - Current User - No Auth after logout",
			Method:         "GET",
			URL:            "{{APIURL}}/user",
			TokenName:      "token1",
			ResponseStatus: 401,
		},
	}

	for _, item := range testCases {
		ok := t.Run(item.Name, func(t *testing.T) {

			if item.Before != nil {
				item.Before()
			}
			// some kind of eval params with substitution
			if item.Expected != nil {
				item.Expected = item.Expected.(func() interface{})()
			}

			var (
				body []byte
				url  = replaceRe.ReplaceAllFunc([]byte(item.URL), replacer)
			)
			if item.Body != "" {
				body = replaceRe.ReplaceAllFunc([]byte(item.Body), replacer)
			}

			req, _ := http.NewRequest(item.Method, string(url), bytes.NewReader(body))
			req.Header.Add("X-Requested-With", "XMLHttpRequest")
			req.Header.Add("Content-Type", "application/json")

			if item.TokenName != "" {
				req.Header.Add("Authorization", "Token "+tplParams[item.TokenName])
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)

			// t.Logf("\nreq body: %s\nresp body: %s", body, respBody)

			if item.ResponseStatus != resp.StatusCode {
				t.Fatalf("bad status code, want: %v, have:%v", item.ResponseStatus, resp.StatusCode)
			}

			// for cases with just status check
			if item.Expected == nil {
				return
			}

			got := WeirdMagicClone(item.Expected)
			err = json.Unmarshal(respBody, got)
			if err != nil {
				t.Fatalf("cant unmarshal resp: %s, body: %s", err, respBody)
			}

			diff, equal := messagediff.PrettyDiff(item.Expected, got)
			if !equal {
				t.Fatalf("\033[1;31mresults not match\033[0m\n \033[1;35mbody\033[0m: %s\n\033[1;32mwant\033[0m %#v\n\033[1;34mgot\033[0m %#v\n\033[1;33mdiff\033[0m:\n%s", respBody, item.Expected, got, diff)
			}

			if item.After != nil {
				err = item.After(resp, respBody, got)
				if err != nil {
					t.Fatalf("after func failed %s", err)
				}
			}
		})
		if !ok {
			break
		}
	}
}

func __test_dummy() {
	fmt.Println(123)
}

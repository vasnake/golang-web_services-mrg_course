package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	ioutil "io" // "io/ioutil"
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

type CR map[string]interface{}

type GQLParams struct {
	Query     string `json:"query"`
	Variables CR     `json:"variables"`
}

type ApiTestCase struct {
	Name           string
	Event          []string
	Method         string
	BodyRaw        string
	GQL            string
	GQLVars        CR
	URL            string
	TokenName      string
	ResponseStatus int
	ResponsePath   string
	Expected       interface{}
	ExpectedRaw    string
	Before         func()
	After          func(*http.Response, []byte, interface{}) error
	CheckFunc      func(interface{}) error
}

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

func WeirdMagicClone(in interface{}) interface{} {
	return reflect.New(reflect.TypeOf(in).Elem()).Interface()
}

func JSONString(in interface{}) string {
	data, _ := json.Marshal(in)
	return string(data)
}

func TestApp(t *testing.T) {
	var (
		app = GetApp()
		ts  = httptest.NewServer(app)

		gqlURL = "/query"

		username = "golang"
		password = "love"
	)

	tplParams := map[string]string{
		"EMAIL":    username + "@example.com",
		"PASSWORD": password,
		"USERNAME": username,
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
			Name: "Catalogs list",
			GQL: `
            {
                Catalog(ID: "1") {
                  id
                  name
                  childs {
                    id
                    name
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 1,
                    "name": "ShopQL",
                    "childs": [
                      {
                        "id": 2,
                        "name": "Книги"
                      },
                      {
                        "id": 5,
                        "name": "Чай"
                      }
                    ]
                  }
                }
              }
            `,
		}, // 0
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog page with param and items list - Tea - with default limit",
			GQL: `
            {
                Catalog(ID: "5") {
                  id
                  name
                  items {
                    id
                    name
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 5,
                    "name": "Чай",
                    "items": [
                      {
                        "id": 9,
                        "name": "Си Пу Юань, Шен Пуэр"
                      },
                      {
                        "id": 10,
                        "name": "Мэнхай 7542, Шен Пуэр"
                      },
                      {
                        "id": 11,
                        "name": "Дянь Хун"
                      }
                    ]
                  }
                }
              }
            `,
		}, // 1
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog page with param and items list - Books - items in subcatalog",
			GQL: `
            {
                Catalog(ID: "2") {
                  id
                  name
                  childs {
                    id
                    name
                    items {
                      id
                      name
                    }
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 2,
                    "name": "Книги",
                    "childs": [
                      {
                        "id": 3,
                        "name": "Алгоритмы",
                        "items": [
                          {
                            "id": 1,
                            "name": "Грокаем алгоритмы | Бхаргава Адитья"
                          },
                          {
                            "id": 2,
                            "name": "Теоретический минимум по Computer Science | Фило Владстон Феррейра"
                          },
                          {
                            "id": 3,
                            "name": "Совершенный алгоритм. Основы | Рафгарден Тим"
                          }
                        ]
                      },
                      {
                        "id": 4,
                        "name": "Golang",
                        "items": [
                          {
                            "id": 5,
                            "name": "Язык программирования Go | Донован Алан А. А., Керниган Брайан У."
                          },
                          {
                            "id": 6,
                            "name": "Go на практике | Butcher Matt, Фарина Мэтт Мэтт"
                          },
                          {
                            "id": 7,
                            "name": "Программирование на Go. Разработка приложений XXI века | Саммерфильд Марк"
                          }
                        ]
                      }
                    ]
                  }
                }
              }
            `,
		}, // 2
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog page with param and items list and pagination",
			GQL: `
            {
                Catalog(ID: "5") {
                    id
                    name
                    items(limit: 1, offset: 1) {
                    id
                    name
                    }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                    "Catalog": {
                    "id": 5,
                    "name": "Чай",
                    "items": [
                        {
                        "id": 10,
                        "name": "Мэнхай 7542, Шен Пуэр"
                        }
                    ]
                    }
                }
                }
            `,
		}, // 3
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog with seller name",
			GQL: `
            {
                Catalog(ID: "5") {
                  id
                  name
                  items(limit: 1, offset: 1) {
                    id
                    name
                    seller {
                      name
                    }
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 5,
                    "name": "Чай",
                    "items": [
                      {
                        "id": 10,
                        "name": "Мэнхай 7542, Шен Пуэр",
                        "seller": {
                          "name": "Дядюшка Ляо"
                        }
                      }
                    ]
                  }
                }
              }
            `,
		}, // 4
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog with other seller",
			GQL: `
            {
                Catalog(ID: "3") {
                  id
                  name
                  items(limit: 5) {
                    id
                    name
                    seller {
                      name
                    }
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 3,
                    "name": "Алгоритмы",
                    "items": [
                      {
                        "id": 1,
                        "name": "Грокаем алгоритмы | Бхаргава Адитья",
                        "seller": {
                          "name": "Издательство Питер"
                        }
                      },
                      {
                        "id": 2,
                        "name": "Теоретический минимум по Computer Science | Фило Владстон Феррейра",
                        "seller": {
                          "name": "Издательство Питер"
                        }
                      },
                      {
                        "id": 3,
                        "name": "Совершенный алгоритм. Основы | Рафгарден Тим",
                        "seller": {
                          "name": "Издательство Питер"
                        }
                      },
                      {
                        "id": 4,
                        "name": "Алгоритмы на Java | Джитер Кевин Уэйн, Седжвик Роберт",
                        "seller": {
                          "name": "Издательство Вильямс"
                        }
                      }
                    ]
                  }
                }
            }
            `,
		}, // 5
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog - inStockText",
			GQL: `
            {
                Catalog(ID: "5") {
                  id
                  name
                  items(limit: 5) {
                    id
                    name
                    inStockText
                  }
                }
            }			  
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 5,
                    "name": "Чай",
                    "items": [
                      {
                        "id": 9,
                        "name": "Си Пу Юань, Шен Пуэр",
                        "inStockText": "мало"
                      },
                      {
                        "id": 10,
                        "name": "Мэнхай 7542, Шен Пуэр",
                        "inStockText": "хватает"
                      },
                      {
                        "id": 11,
                        "name": "Дянь Хун",
                        "inStockText": "хватает"
                      },
                      {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inStockText": "много"
                      },
                      {
                        "id": 13,
                        "name": "Габа Улун",
                        "inStockText": "много"
                      }
                    ]
                  }
                }
              }
            `,
		}, // 6
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Seller with items and item catalog",
			GQL: `
            {
                Seller(ID: "3") {
                  id
                  name
                  items(limit: 5) {
                    id
                    name
                    parent {
                      id
                      name
                    }
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "data": {
                  "Seller": {
                    "id": 3,
                    "name": "Издательство Питер",
                    "items": [
                      {
                        "id": 1,
                        "name": "Грокаем алгоритмы | Бхаргава Адитья",
                        "parent": {
                          "id": 3,
                          "name": "Алгоритмы"
                        }
                      },
                      {
                        "id": 2,
                        "name": "Теоретический минимум по Computer Science | Фило Владстон Феррейра",
                        "parent": {
                          "id": 3,
                          "name": "Алгоритмы"
                        }
                      },
                      {
                        "id": 3,
                        "name": "Совершенный алгоритм. Основы | Рафгарден Тим",
                        "parent": {
                          "id": 3,
                          "name": "Алгоритмы"
                        }
                      },
                      {
                        "id": 8,
                        "name": "Head First. Изучаем Go | Макгаврен Джей",
                        "parent": {
                          "id": 4,
                          "name": "Golang"
                        }
                      }
                    ]
                  }
                }
              }
            `,
		}, // 7
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog - how many in cart - ERROR(no access) - directive @authorized",
			GQL: `
            {
                Catalog(ID: "5") {
                  id
                  name
                  items(limit: 1) {
                    id
                    name
                    inCart
                  }
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "errors": [
                  {
                    "message": "User not authorized",
                    "path": [
                      "Catalog",
                      "items",
                      0,
                      "inCart"
                    ]
                  }
                ],
                "data": {
                  "Catalog": null
                }
            }
            `,
		}, // 8
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Add to cart - ERROR(no access) - directive @authorized",
			GQL: `
            mutation {
                AddToCart(in: {itemID: 1, quantity: 2}) {
                  item {
                    id
                    name
                  }
                  quantity
                }
            }
            `,
			URL: gqlURL,
			ExpectedRaw: `
            {
                "errors": [
                  {
                    "message": "User not authorized",
                    "path": [
                      "AddToCart"
                    ]
                  }
                ],
                "data": null
            }
            `,
		}, // 9
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name:           "Register",
			URL:            "/register",
			Method:         http.MethodPost,
			BodyRaw:        "{\"user\":{\"email\":\"{{EMAIL}}\", \"password\":\"{{PASSWORD}}\", \"username\":\"{{USERNAME}}\"}}",
			ResponseStatus: 200,
			CheckFunc: func(resp interface{}) error {
				val, err := lookup.LookupString(resp, "body.token")
				if err != nil {
					return err
				}
				fmt.Println("-------------------- TOKEN:", val)
				tplParams["token1"] = val.String()
				return nil
			},
		}, // 10
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Add to cart - first item - success",
			GQL: `
            mutation {
                AddToCart(in: {itemID: 12, quantity: 2}) {
                  item {
                    id
                    name
                    inStockText
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "AddToCart": [
                    {
                      "item": {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inStockText": "хватает"
                      },
                      "quantity": 2
                    }
                  ]
                }
              }			`,
		}, // 11
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Add to cart - first item - check correct increment in cart",
			GQL: `
            mutation {
                AddToCart(in: {itemID: 12, quantity: 2}) {
                  item {
                    id
                    name
                    inStockText
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "AddToCart": [
                    {
                      "item": {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inStockText": "мало"
                      },
                      "quantity": 4
                    }
                  ]
                }
            }
            `,
		}, // 12
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Add to cart - first item - check quantity availability",
			GQL: `
            mutation {
                AddToCart(in: {itemID: 12, quantity: 2}) {
                  item {
                    id
                    name
                    inStockText
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "errors": [
                  {
                    "message": "not enough quantity",
                    "path": [
                      "AddToCart"
                    ]
                  }
                ],
                "data": null
            }
            `,
		}, // 13
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Add to cart - second item - before delete check",
			GQL: `
            mutation {
                AddToCart(in: {itemID: 1, quantity: 1}) {
                  item {
                    id
                    name
                    inStockText
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "AddToCart": [
                    {
                      "item": {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inStockText": "мало"
                      },
                      "quantity": 4
                    },
                    {
                        "item": {
                          "id": 1,
                          "name": "Грокаем алгоритмы | Бхаргава Адитья",
                          "inStockText": "мало"
                        },
                        "quantity": 1
                    }
                  ]
                }
            }
            `,
		}, // 14
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Remove from cart",
			GQL: `
            mutation {
                RemoveFromCart(in: {itemID: 1, quantity: 1}) {
                  item {
                    id
                    name
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "RemoveFromCart": [
                    {
                      "item": {
                        "id": 12,
                        "name": "Да Хун Пао"
                      },
                      "quantity": 4
                    }
                  ]
                }
            }
            `,
		}, // 15

		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "My Cart",
			GQL: `
            {
                MyCart {
                  item {
                    id
                    name
                    inStockText
                  }
                  quantity
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "MyCart": [
                    {
                      "item": {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inStockText": "мало"
                      },
                      "quantity": 4
                    }
                  ]
                }
            }
            `,
		}, // 16
		// ----------------------------------------------------------------------------------------
		&ApiTestCase{
			Name: "Catalog page with inCart param",
			GQL: `
            query{
                Catalog(ID: "5") {
                  id
                  name
                  items(limit:8) {
                    id
                    name
                    inCart
                  }
                }
            }
            `,
			URL:       gqlURL,
			TokenName: "token1",
			ExpectedRaw: `
            {
                "data": {
                  "Catalog": {
                    "id": 5,
                    "name": "Чай",
                    "items": [
                      {
                        "id": 9,
                        "name": "Си Пу Юань, Шен Пуэр",
                        "inCart": 0
                      },
                      {
                        "id": 10,
                        "name": "Мэнхай 7542, Шен Пуэр",
                        "inCart": 0
                      },
                      {
                        "id": 11,
                        "name": "Дянь Хун",
                        "inCart": 0
                      },
                      {
                        "id": 12,
                        "name": "Да Хун Пао",
                        "inCart": 4
                      },
                      {
                        "id": 13,
                        "name": "Габа Улун",
                        "inCart": 0
                      }
                    ]
                  }
                }
              }
            `,
		}, // 17
	}

	for _, item := range testCases {
		ok := t.Run(item.Name, func(t *testing.T) {
			if item.Before != nil {
				item.Before()
			}
			// some kind of eval params with substitution
			if item.Expected != nil {
				item.Expected = item.Expected.(func() interface{})()
			} else if item.ExpectedRaw != "" {
				// var data CR
				err := json.Unmarshal([]byte(item.ExpectedRaw), &item.Expected)
				if err != nil {
					t.Fatalf("cant unmarshal json: %v", err)
				}
			}

			var (
				body []byte
				url  = replaceRe.ReplaceAllFunc([]byte(ts.URL+item.URL), replacer)
			)
			if item.GQL != "" {
				item.BodyRaw = JSONString(&GQLParams{
					Query:     item.GQL,
					Variables: item.GQLVars,
				})
			}
			if item.BodyRaw != "" {
				body = replaceRe.ReplaceAllFunc([]byte(item.BodyRaw), replacer)
			}

			// t.Log("body", item.BodyRaw)
			if item.URL == gqlURL {
				item.Method = "POST"
			}
			req, _ := http.NewRequest(item.Method, string(url), bytes.NewReader(body))
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

			// t.Log((item.ResponseStatus == 0 && resp.StatusCode != 200), item.ResponseStatus == 0, resp.StatusCode != 200, item.ResponseStatus, resp.StatusCode, resp.StatusCode == 200)

			if item.ResponseStatus != 0 && (item.ResponseStatus != resp.StatusCode) {
				t.Fatalf("bad status code, want: %v, have:%v", item.ResponseStatus, resp.StatusCode)
			}

			// for cases with just status check
			if item.Expected == nil && item.CheckFunc == nil {
				return
			}
			var got interface{}
			err = json.Unmarshal(respBody, &got)
			if err != nil {
				t.Fatalf("cant unmarshal resp: %s, body: %s", err, respBody)
			}

			// for custom checking logic
			// i'm to lazy to code entire registrtion flow, so it's just check and set token inside
			if item.CheckFunc != nil {
				if err := item.CheckFunc(got); err != nil {
					t.Fatal("CheckFunc failed:", err)
				}
				return
			}

			diff, equal := messagediff.PrettyDiff(item.Expected, got)
			if !equal {
				// dd(item.Expected, got)
				t.Fatalf("\033[1;31mresults not match\033[0m\n\033[1;35mbody\033[0m: %s\n\033[1;32mwant\033[0m %#v\n\033[1;34mgot\033[0m %#v\n\033[1;33mdiff\033[0m:\n%s", respBody, item.Expected, got, diff)
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

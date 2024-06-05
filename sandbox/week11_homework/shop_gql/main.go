package main

import (
	"net/http"

	gql_handler "github.com/99designs/gqlgen/graphql/handler"
	gql_handler_extension "github.com/99designs/gqlgen/graphql/handler/extension"
)

const (
	TEST_DATA_FILE_NAME = "testdata.json"
)

func main() {
	panic("not yet")
}

func GetApp() http.Handler {
	cfg := Config{
		Resolvers: &(Resolver{
			dataAdapter: StorageGQLAdapter{
				shopStorage: loadData(),
			},
		}),
	}

	// https://gqlgen.com/reference/directives/
	cfg.Directives.Authorized = CheckAuthorizedMiddleware

	var gqlSrv = gql_handler.NewDefaultServer(NewExecutableSchema(cfg))
	gqlSrv.Use(gql_handler_extension.FixedComplexityLimit(500))

	mux := http.NewServeMux()
	mux.Handle("/query", gqlSrv)

	var user_sessions_authHandlers = (&UserSessionAuth{}).New()
	mux.HandleFunc("/register", user_sessions_authHandlers.RegisterNewUserHandler)

	handler := user_sessions_authHandlers.InjectSession2ContextMiddleware(mux)

	return handler
}

func loadData() ShopStorage {
	data, err := loadTestData(TEST_DATA_FILE_NAME)
	panicOnError("loadTestData failed", err)

	var sellers []SellerStruct
	sellers, err = loadSellers(data)
	panicOnError("loadSellers failed", err)

	var catalogs []CatalogStruct
	var items []GoodiesItemStruct
	catalogs, items, err = loadCatalogTree(data)
	panicOnError("loadCatalogTree failed", err)

	show("loaded sellers: ", sellers)
	show("loaded items: ", items)
	show("loaded catalogs: ", catalogs)

	return ShopStorage{
		sellersRows: sellers,
		itemsRows:   items,
		catalogRows: catalogs,
	}
}

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

	var srv = gql_handler.NewDefaultServer(NewExecutableSchema(cfg))
	srv.Use(gql_handler_extension.FixedComplexityLimit(500))

	return srv
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

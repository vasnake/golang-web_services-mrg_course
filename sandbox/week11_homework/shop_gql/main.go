package main

import (
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	gqlgen_extension "github.com/99designs/gqlgen/graphql/handler/extension"
)

const (
	TEST_DATA_FILE_NAME = "testdata.json"
)

func main() {
	panic("not yet")
}

func GetApp() http.Handler {
	data, err := loadTestData(TEST_DATA_FILE_NAME)
	panicOnError("loadTestData failed", err)
	// show("loaded data: ", data)

	var sellers []SellerStruct
	sellers, err = loadSellers(data)
	panicOnError("loadSellers failed", err)
	show("loaded sellers: ", sellers)
	var catalogs []CatalogStruct
	var items []GoodiesItemStruct
	catalogs, items, err = loadCatalogTree(data)
	panicOnError("loadCatalogTree failed", err)
	show("loaded catalogs: ", catalogs)
	show("loaded items: ", items)

	var storage = ShopStorage{
		sellersRows: sellers,
		itemsRows:   items,
		catalogRows: catalogs,
	}
	show("storage: ", storage)
	// var sga = StorageGQLAdapter{shopStorage: storage}
	// err = sga.rebuildLists()
	// panicOnError("adapter rebuildLists failed", err)
	// gqlResolver := &Resolver{dataAdapter: sga}
	gqlResolver := &Resolver{}

	cfg := Config{
		Resolvers: gqlResolver,
	}

	srv := graphql_handler.NewDefaultServer(NewExecutableSchema(cfg))
	srv.Use(gqlgen_extension.FixedComplexityLimit(500))

	return srv
}

type ShopStorage struct {
	sellersRows []SellerStruct
	itemsRows   []GoodiesItemStruct
	catalogRows []CatalogStruct
}

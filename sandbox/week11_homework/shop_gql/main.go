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
	// state
	r := (&Resolver{}).New()
	r.dataAdapter.shopStorage = loadData()

	// graphql service
	gqlCfg := Config{Resolvers: r}
	gqlCfg.Directives.Authorized = CheckAuthorizedMiddleware // https://gqlgen.com/reference/directives/
	var gqlSvc = gql_handler.NewDefaultServer(NewExecutableSchema(gqlCfg))
	gqlSvc.Use(gql_handler_extension.FixedComplexityLimit(500))

	routesMux := http.NewServeMux()

	// route for graphql handlers
	routesMux.Handle("/query", gqlSvc)
	// route for user registration handlers
	var authSvc = (&UserSessionAuthSvc{}).New()
	routesMux.HandleFunc("/register", authSvc.RegisterNewUserHandler)

	// add middleware
	app := authSvc.InjectSession2ContextMiddleware(routesMux)

	return app
}

func loadData() *ShopStorage {
	data, err := loadTestData(TEST_DATA_FILE_NAME)
	panicOnError("loadTestData failed", err)

	var sellers []*SellerStruct
	var catalogs []*CatalogStruct
	var items []*GoodiesItemStruct

	sellers, err = loadSellers(data)
	panicOnError("loadSellers failed", err)

	catalogs, items, err = loadCatalogTree(data)
	panicOnError("loadCatalogTree failed", err)

	show("loaded sellers: ", sellers)
	show("loaded items: ", items)
	show("loaded catalogs: ", catalogs)

	return &ShopStorage{
		sellersRows: sellers,
		itemsRows:   items,
		catalogRows: catalogs,
	}
}

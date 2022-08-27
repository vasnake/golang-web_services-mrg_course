package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// handler, N.B. signature
func Index(ctx *fasthttp.RequestCtx) {
	users := []string{"rvasily"}
	body, _ := json.Marshal(users)

	ctx.SetBody(body)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "you try to see user %s\n", ctx.UserValue("id"))
}

func main() {
	router := fasthttprouter.New()

	router.GET("/", Index)
	router.GET("/users/:id", GetUser)

	fmt.Println("starting server at :8080")
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

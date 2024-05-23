package oauth_demo

import (
	"encoding/json"
	"flag"
	"fmt"
	ioutil "io" // "io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

const (
	/*
		manage app: https://github.com/settings/applications/2587428
		app name: go-webservice-course-vromanov
		homepage: http://localhost:8080/
		callback: http://localhost:8080/user/login_oauth
		Client ID: Ov***gJF
		Client secret: ada***860
		docs: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authenticating-to-the-rest-api-with-an-oauth-app
	*/
	REDIRECT_URL = "http://localhost:8080/user/login_oauth" // callback
	AUTH_URL     = "https://github.com/login/oauth/authorize?scope=user:email&client_id=%s"
	API_URL      = "https://api.github.com/user?fields=email,photo_50&access_token=%s"
)

var (
	// from command line, using flag package
	// export OAUTH_APP_ID=Ov2***gJF
	// export OAUTH_APP_SECRET=ada***860
	// go run week10 -appid ${OAUTH_APP_ID:-foo} -appsecret ${OAUTH_APP_SECRET:-bar}
	APP_ID     = "Ov2***gJF"
	APP_SECRET = "ada***860"
)

func MainDemo() {
	flag.StringVar(&APP_ID, "appid", "foo?", "app id (client id) from github registered app")
	flag.StringVar(&APP_SECRET, "appsecret", "bar?", "app secret (client key) from github registered app")
	flag.Parse()
	show("you mustn't but: appid, appsecret: ", APP_ID, APP_SECRET)

	http.HandleFunc("/", rootHandler)
	log.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// actually, here we implemented two different handlers:
	// 1) login handler, where user will be redirected to provider auth endpoint;
	// 2) callback aka redirect url, where user will be redirected from provider, after auth dance

	// after using AUTH_URL by user, provider give as a code (on step 2)
	code := r.FormValue("code") // code from provider, e.g. https://example.com/user/login?code=123456

	// no code yet, let's start with AUTH_URL (handler 1)
	if code == "" {
		var href = fmt.Sprintf(AUTH_URL, APP_ID)
		w.Write([]byte(`<div><a href="` + href + `">authorize</a></div>`))
		return // wait for client to decide what to do
	}

	// handler 2
	// after executing auth dance user-provider, provider ask user to go to 'redirect url' aka callback,
	// and user teleports here: in our app again, but we have code now,
	// and we need token from provider in exchange for that code

	conf := oauth2.Config{
		ClientID:     APP_ID,
		ClientSecret: APP_SECRET,
		RedirectURL:  REDIRECT_URL,
		// Endpoint:     vk.Endpoint,
		Endpoint: GitHubEndpoint,
	}
	ctx := r.Context()
	oauthToken, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println("cannot exchange code to token", err)
		http.Error(w, err.Error(), 500)
		return
	}
	show("got oauth token: ", oauthToken)
	// got token (just token), no email or user_id
	// &oauth2.Token{AccessToken:"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", TokenType:"bearer", RefreshToken:"", Expiry:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), raw:url.Values{"access_token":[]string{"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"}, "scope":[]string{"user:email"}, "token_type":[]string{"bearer"}}, expiryDelta:0}

	// nah, nothing here:
	email := oauthToken.Extra("email").(string)
	show("got email from oauth token: ", email)
	// userID_float := oauthToken.Extra("user_id").(float64) // WTF?
	userID := oauthToken.Extra("user_id").(string)
	show("got user id from oauth token: ", userID)
	// userID := int(userID_float)

	// show token to user

	w.Write([]byte(`
	<div> Oauth token (provider provided in exchange to code):<br>
		` + fmt.Sprintf("%#v", oauthToken) + `
	</div> <br>
	<div>Email (from token): '` + email + `'</div> <br>
	<div>UserID (from token): '` + userID + `'</div>
	<br>
	`))

	// we want to ask provider for some user data
	httpClient := conf.Client(ctx, oauthToken)
	apiResp, err := httpClient.Get(fmt.Sprintf(API_URL, oauthToken.AccessToken))
	if err != nil {
		log.Println("cannot request data from provider (api or token not working well)", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer apiResp.Body.Close()

	// decode api response
	respBodyBytes, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		log.Println("cannot read buffer", err)
		http.Error(w, err.Error(), 500)
		return
	}
	show("api response: ", string(respBodyBytes))
	userData := make(map[string]any, 32)
	err = json.Unmarshal(respBodyBytes, &userData)
	if err != nil {
		log.Println("cannot json.Unmarshal", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if len(userData) == 0 {
		log.Println("requested data is empty", err)
		http.Error(w, "you should read the api docs", 500)
		return
	}

	// extract some data
	emailAny, emailExists := userData["email"]
	if emailExists {
		email = emailAny.(string)
	}
	uidAny, uidExists := userData["id"]
	if uidExists {
		userID = strconv.FormatUint(uint64(uidAny.(float64)), 36)
		// userID = strconv.FormatUint(uint64(uidAny.(float64)), 10)
	}
	nameAny, _ := userData["name"] // nah, error handling: go fuck yourself
	avatarAny, _ := userData["avatar_url"]

	// show data to user
	w.Write([]byte(`
	<div>
		<img src="` + avatarAny.(string) + `"/> <br>
		` + nameAny.(string) + `
	</div>
	<br>
	<div>Email (from api): '` + email + `'</div> <br>
	<div>UserID (from api): '` + userID + `'</div>
	<br>	
	`))
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}

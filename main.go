package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

const (
	AuthorizeURL = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"
)

type oAuthResponse struct {
	Code         string `json:"code"`
	AuthorizeURL string `json:"authorize_url"`
}

type authorizeResponse struct {
	AccessToken string `json:"access_token"`
}

type getAPIResponse struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	ItemType string `json:"item_type"`
}

func main() {
	http.HandleFunc("/oauth/request", func(w http.ResponseWriter, r *http.Request) {
		urls, ok := r.URL.Query()["redirect_url"]

		if !ok || len(urls[0]) < 1 {
			http.Error(w, "Url Param 'url' is missing", http.StatusBadRequest)
			return
		}

		code, err := OAuthRequest(urls[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := new(oAuthResponse)
		data.AuthorizeURL = fmt.Sprintf(AuthorizeURL, code, urls[0])
		data.Code = code
		res, err := json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		urls, ok := r.URL.Query()["code"]

		if !ok || len(urls[0]) < 1 {
			http.Error(w, "Url Param 'url' is missing", http.StatusBadRequest)
			return
		}

		token, err := AuthorizeRequest(urls[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := new(authorizeResponse)
		data.AccessToken = token

		res, err := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.URL.Query()["access_token"]

		if !ok || len(tokens[0]) < 1 {
			http.Error(w, "Url Param 'tokens' is missing", http.StatusBadRequest)
			return
		}

		apiRes, err := GetRequest(tokens[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]getAPIResponse, 0)

		for _, v := range apiRes.List {
			items = append(items, getAPIResponse{Title: v.GivenTitle, Url: v.GivenURL, ItemType: "pocket"})
		}

		res, err := json.Marshal(items)

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"),
		handlers.LoggingHandler(os.Stdout,
			LimitHandler(http.DefaultServeMux))); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

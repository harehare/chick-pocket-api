package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	OAuthAPIURL     = "https://getpocket.com/v3/oauth/request"
	AuthorizeAPIURL = "https://getpocket.com/v3/oauth/authorize"
	GetAPIURL       = "https://getpocket.com/v3/get"
)

type RequestToken struct {
	Code string `json:"code"`
}

type Authorization struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

type Item struct {
	ItemID        string `json:"item_id"`
	ResolvedID    string `json:"resolved_id"`
	GivenURL      string `json:"given_url"`
	GivenTitle    string `json:"given_title"`
	Favorite      string `json:"favorite"`
	Status        string `json:"status"`
	ResolvedTitle string `json:"resolved_title"`
	ResolvedURL   string `json:"resolved_url"`
	Excerpt       string `json:"excerpt"`
	IsArticle     string `json:"is_article"`
	HasVideo      string `json:"has_video"`
	HasImage      string `json:"has_image"`
	WordCount     string `json:"word_count"`
	Images        struct {
		Image struct {
			ItemID  string `json:"item_id"`
			ImageID string `json:"image_id"`
			Src     string `json:"src"`
			Width   string `json:"width"`
			Height  string `json:"height"`
			Credit  string `json:"credit"`
			Caption string `json:"caption"`
		} `json:"image"`
	} `json:"images"`
	Videos struct {
		Video struct {
			ItemID  string `json:"item_id"`
			VideoID string `json:"video_id"`
			Src     string `json:"src"`
			Width   string `json:"width"`
			Height  string `json:"height"`
			Type    string `json:"type"`
			Vid     string `json:"vid"`
		} `json:"video"`
	} `json:"videos"`
}

type Retrieve struct {
	Status   int             `json:"status"`
	Complete int             `json:"complete"`
	List     map[string]Item `json:"list"`
}

func postJSON(url string, jsonData map[string]string) ([]byte, error) {
	body, err := json.Marshal(jsonData)
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetRequest(accessToken string) (*Retrieve, error) {

	body, err := postJSON(
		GetAPIURL,
		map[string]string{
			"consumer_key": os.Getenv("CONSUMER_KEY"),
			"access_token": accessToken,
			"sort":         "newest",
			"count":        "100",
			"detailType":   "complete",
		})

	if err != nil {
		return nil, err
	}

	data := new(Retrieve)

	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	return data, nil
}

func AuthorizeRequest(code string) (string, error) {

	body, err := postJSON(
		AuthorizeAPIURL,
		map[string]string{
			"consumer_key": os.Getenv("CONSUMER_KEY"),
			"code":         code,
		})

	if err != nil {
		return "", err
	}

	data := new(Authorization)

	if err := json.Unmarshal(body, data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}

func OAuthRequest(redirectURI string) (string, error) {

	body, err := postJSON(OAuthAPIURL, map[string]string{
		"consumer_key": os.Getenv("CONSUMER_KEY"),
		"redirect_uri": redirectURI,
	})

	if err != nil {
		return "", err
	}

	data := new(RequestToken)

	if err := json.Unmarshal(body, data); err != nil {
		return "", err
	}

	return data.Code, nil
}

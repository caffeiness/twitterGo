package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	// Json読み込み
	raw, error := ioutil.ReadFile("twitterAccount.json")
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	var twitterAccount TwitterAccount
	// 構造体にセット
	json.Unmarshal(raw, &twitterAccount)

	// 認証
	api := anaconda.NewTwitterApiWithCredentials(twitterAccount.AccessToken, twitterAccount.AccessTokenSecret, twitterAccount.ConsumerKey, twitterAccount.ConsumerSecret)

	v := url.Values{}
	v.Set("count", "30")
	tweets, err := api.GetHomeTimeline(v)
	if err != nil {
		panic(err)
	}

	for _, tweet := range tweets {
		fmt.Println("tweet: ", tweet.Text)
	}
}

// TwitterAccount はTwitterの認証用の情報
type TwitterAccount struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
}

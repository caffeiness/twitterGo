package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
)

func connectTwitterApi() *anaconda.TwitterApi {
	//Key情報のあるjson読み込み
	raw, error := ioutil.ReadFile("twitterAccount.json")
	if error != nil {
		fmt.Println(error.Error())
		return nil
	}

	var twitterAccount TwitterAccount
	json.Unmarshal(raw, &twitterAccount)

	// 認証
	return anaconda.NewTwitterApiWithCredentials(twitterAccount.AccessToken, twitterAccount.AccessTokenSecret, twitterAccount.ConsumerKey, twitterAccount.ConsumerSecret)
}

func serach(c echo.Context) error {
	keyword := c.FormValue("keyword")
	api := connectTwitterApi()
	// 検索
	searchResult, _ := api.GetSearch(`"`+keyword+`"`, nil)

	tweets := make([]*Tweet, 0)

	for _, data := range searchResult.Statuses {
		tweet := new(Tweet)
		tweet.Text = data.FullText
		tweet.User = data.User.Name

		tweets = append(tweets, tweet)
	}

	return c.JSON(http.StatusOK, tweets)
}

func main() {
	e := echo.New()
	e.POST("/tweet", serach)
	e.Logger.Fatal(e.Start(":1323"))
}

type TwitterAccount struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
}

type Tweet struct {
	User string `json:"user"`
	Text string `json:"text"`
}

type Tweets *[]Tweet

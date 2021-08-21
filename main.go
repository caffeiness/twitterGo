package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

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

func tweets(c echo.Context) error {
	value := c.QueryParam("value")
	api := connectTwitterApi()
	//検索
	searchResult, _ := api.GetSearch(`"`+value+`"`, nil)
	tweets := make([]*TweetTempete, 0)
	for _, data := range searchResult.Statuses {
		tweet := new(TweetTempete)
		tweet.Text = data.FullText
		tweet.User = data.User.Name
		tweet.Id = data.User.IdStr
		tweet.ScreenName = data.User.ScreenName
		tweet.Date = data.CreatedAt
		tweet.TweetId = data.IdStr
		tweets = append(tweets, tweet)
	}

	return c.Render(http.StatusOK, "index.html", tweets)
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("index.html")),
	}
	e.Renderer = t
	e.GET("/tweet", tweets)
	e.Logger.Fatal(e.Start(":1323"))
}

type TwitterAccount struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
}

type Template struct {
	templates *template.Template
}

// TweetTempete はツイートの情報
type TweetTempete struct {
	User       string `json:"user"`
	Text       string `json:"text"`
	ScreenName string `json:"screenName"`
	Id         string `json:"id"`
	Date       string `json:"date"`
	TweetId    string `json:"tweetId"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

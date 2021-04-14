package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	//redditID       string
	//redditSecret   string
	//redditUsername string
	//redditPassword string
	TelegramChatID  string   `envconfig:"TELEGRAM_CHAT_ID"`
	TelegramToken   string   `envconfig:"TELEGRAM_TOKEN"`
	SearchSubreddit string   `envconfig:"REDDIT_SEARCH_SUBREDDIT"`
	SearchTerms     []string `envconfig:"REDDIT_SEARCH_TERMS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	//spew.Dump(config)

	//credentials := reddit.Credentials{ID: config.redditID, Secret: config.redditSecret, Username: config.redditUsername, Password: config.redditPassword}
	//client, err := reddit.NewClient(http.DefaultClient, &credentials)
	//if err != nil {
	//	log.Fatal(err)
	//}
	client := reddit.DefaultClient()

	for {
		log.Info(fmt.Sprintf("about to search for the following terms: [%s] in r/buildapcsales", strings.Join(config.SearchTerms, ",")))
		err := searchAndNotify(client, &config)
		if err != nil {
			log.Error("main: something went wrong", err)
		}
		log.Info("will check again in 5 min!")
		time.Sleep(5 * time.Minute)
	}
}

var latestResultID = make(map[string]string)
var searchOptions = reddit.ListPostSearchOptions{
	ListPostOptions: reddit.ListPostOptions{
		Time: "day",
	},
	Sort: "new",
}

func searchAndNotify(c *reddit.Client, config *Config) error {
	for _, searchTerm := range config.SearchTerms {
		so := searchOptions
		// if we've already searched this term before, only include new results
		if id, ok := latestResultID[searchTerm]; ok {
			so.After = id
		}

		posts, _, err := c.Subreddit.SearchPosts(context.Background(), searchTerm, config.SearchSubreddit, &so)
		if err != nil {
			log.Error("searchAndNotify: could not reddit posts", err)
			return err
		}
		numPosts := len(posts)
		log.Info(fmt.Sprintf("found %d posts matching `%s`!", numPosts, searchTerm))

		if numPosts <= 0 {
			break
		}
		// Save this as an 'anchor' point to remember the last found result
		latestResultID[searchTerm] = posts[0].FullID

		for _, p := range posts {
			if p == nil { // this shouldn't happen but who knows anymore ¯\_(ツ)_/¯
				continue
			}
			err := sendTelegramMessage(config, p.Title, p.URL)
			if err != nil {
				log.Error("searchAndNotify: could not send message", err)
				return err
			}
		}
	}

	return nil
}

var telegramRequestURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s"

func sendTelegramMessage(config *Config, title, link string) error {
	msgBody := `
%s

%s`
	msgBody = fmt.Sprintf(msgBody, title, link)
	u := fmt.Sprintf(telegramRequestURL, config.TelegramToken, config.TelegramChatID, url.QueryEscape(msgBody))
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		log.Error("sendTelegramMessage: could build proper request")
		return err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("sendTelegramMessage: could not send telegram message")
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode >= http.StatusOK && rsp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	log.Error("sendTelegramMessage: something went wrong with the request")
	body, err := ioutil.ReadAll(rsp.Body)
	log.Info("body: ", string(body))

	return errors.New("sendTelegramMessage: could not send telegram message")
}

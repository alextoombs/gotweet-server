package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ChimeraCoder/anaconda"
)

// twitterApi is the initial API created from access tokens found on disk or in
// the environment.
var twitterApi *anaconda.TwitterApi

func main() {
	// Set up a port to listen on.
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	accessTokens, err := getTokens()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Flush state to disk.
	if err := accessTokens.save(); err != nil {
		log.Fatal(err.Error())
	}

	twitterApi = createTwitterApi(accessTokens)

	http.HandleFunc("/tweet", tweetPostHandler)
	log.Printf("Serving traffic on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error service traffic: ", err)
	}
}

// createTwitterApi creates a new Twitter API object from a set of tokens.
func createTwitterApi(t *tokens) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	return anaconda.NewTwitterApi(t.AccessToken, t.AccessTokenSecret)
}

// tweetPostHandler posts the string body of the request to twitter.
func tweetPostHandler(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer req.Body.Close()

	log.Printf("Received request body to POST: %s\n", string(b))
	tweet, err := twitterApi.PostTweet(string(b), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Posted Tweet: %v\n", tweet)
}

const settingsFile = ".gotweet-server"

// gotweetPath is the full path to where the gotweet-server settings file lives.
var gotweetPath = filepath.Join(os.Getenv("HOME"), settingsFile)

// getTokens will first attempt to retrieve tokens from a file located at
// gotweetPath. If not found there, it will attempt to retrieve tokens from
// environment. If that fails, then it will throw an error, aborting the startup
// of the server. tokens structs returned are guaranteed to be validated.
func getTokens() (*tokens, error) {
	f, err := os.Open(gotweetPath)
	if err != nil {
		// If the file doesn't exist, that's okay. We attempt to retrieve tokens
		// from the environment. If those aren't all valid, then we error.
		if os.IsNotExist(err) {
			t := getTokensFromEnvironment()
			if !t.isValid() {
				return nil, fmt.Errorf("could not retrieve tokens from environment")
			}
			return t, nil
		}
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var t *tokens
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	if !t.isValid() {
		return nil, fmt.Errorf("could not retrieve all tokens from disk")
	}
	return t, nil
}

// getTokensFromEnvironment populates all necessary tokens for Twitter operation
// from the environment.
func getTokensFromEnvironment() *tokens {
	return &tokens{
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	}
}

// tokens is a collection of all the important tokens needed to use a twitter
// client. This may be written to disk in order to remove secrets from the
// environment.
type tokens struct {
	ConsumerKey       string `json:"consumer_key"`    // Maps to API key
	ConsumerSecret    string `json:"consumer_secret"` // Maps to API secret
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

// isValid returns true if all necessary keys are populated.
func (t tokens) isValid() bool {
	return !(t.ConsumerKey == "" || t.ConsumerSecret == "" ||
		t.AccessToken == "" || t.AccessTokenSecret == "")
}

// save writes a set of tokens to disk.
func (t tokens) save() error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	tmpPath := filepath.Join(os.TempDir(), "temp_gotweet")
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	// Move into correct path.
	return os.Rename(tmpPath, gotweetPath)
}

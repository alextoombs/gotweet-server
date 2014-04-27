package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetTokens(t *testing.T) {
	// No env, no file.
	toks, err := getTokens()
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != "could not retrieve tokens from environment" {
		t.Fatalf("Expected different error, not: %s", err.Error())
	}

	mockEnvironmentTokens(t)
	toks, err = getTokens()
	if err != nil || toks == nil {
		t.Fatalf("Expected no error, not: %s", err.Error())
	}
}

func TestGetTokens_InvalidFile(t *testing.T) {
	mockTokensFile_Corrupted(t)

	// No env, no file.
	toks, err := getTokens()
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != "could not retrieve all tokens from disk" {
		t.Fatalf("Expected different error, not: %s", err.Error())
	}

	mockTokensFile(t)
	toks, err = getTokens()
	if err != nil || toks == nil {
		t.Fatalf("Expected no error, not: %s", err.Error())
	}
}

func mockEnvironmentTokens(t *testing.T) {
	if err := os.Setenv("TWITTER_CONSUMER_KEY", "foo"); err != nil {
		t.Fatalf("Should be able to set environment: %s", err.Error())
	}
	if err := os.Setenv("TWITTER_CONSUMER_SECRET", "foo"); err != nil {
		t.Fatalf("Should be able to set environment: %s", err.Error())
	}
	if err := os.Setenv("TWITTER_ACCESS_TOKEN", "foo"); err != nil {
		t.Fatalf("Should be able to set environment: %s", err.Error())
	}
	if err := os.Setenv("TWITTER_ACCESS_TOKEN_SECRET", "foo"); err != nil {
		t.Fatalf("Should be able to set environment: %s", err.Error())
	}
}

func mockTokensFile(t *testing.T) {
	gotweetPath = filepath.Join(os.TempDir(), "mockgotweetsettings")
	tokens := &tokens{
		ConsumerKey:       "foo",
		ConsumerSecret:    "bar",
		AccessToken:       "baz",
		AccessTokenSecret: "blah",
	}
	if err := tokens.save(); err != nil {
		t.Fatalf("Could not save tokens: %s", err)
	}
}

func mockTokensFile_Corrupted(t *testing.T) {
	gotweetPath = filepath.Join(os.TempDir(), "mockgotweetsettings")
	tokens := &tokens{
		ConsumerSecret: "bar",
		AccessToken:    "baz",
	}
	if err := tokens.save(); err != nil {
		t.Fatalf("Could not save tokens: %s", err)
	}
}

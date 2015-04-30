package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var githubTokenValue string

func GithubToken() string {
	if githubTokenValue == "" {
		githubTokenValue = githubToken()
	}
	return githubTokenValue
}

func githubToken() string {
	tokenBuf, err := exec.Command("git", "config", "--get", "github.token").CombinedOutput()
	token := strings.TrimSpace(string(tokenBuf))

	if err != nil || token == "" {
		fmt.Printf("You are missing a GitHub authorization token in your git configuration\n")
		fmt.Printf("Get a token from https://github.com/settings/tokens/new\n")
		fmt.Printf("Place it in your git configuration with something like:\n")
		fmt.Printf("git config --set github.token YOURTOKEN\n")
		os.Exit(1)
	}
	return token
}

func GithubApi(method string, uri string, requestBody map[string]interface{}) (map[string]interface{}, error) {
	var requestBodyBuf io.Reader
	if requestBody != nil {
		buf, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		requestBodyBuf = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method,
		"https://api.github.com"+uri, requestBodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+GithubToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 400 {
		return nil, err
	}

	body := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

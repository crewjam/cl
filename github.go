package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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
		fmt.Printf("git config --add github.token YOURTOKEN\n")
		os.Exit(1)
	}
	return token
}

func GithubRepo() string {
	buf, err := exec.Command("git", "config", "--local", "--get", "github.repo").CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(buf))
	}

	buf, err = exec.Command("git", "config", "--local", "--get", "remote.origin.url").CombinedOutput()
	if err == nil {
		matches := regexp.MustCompile(`^git@github.com:(.*).git\n$`).FindAllStringSubmatch(string(buf), -1)
		if matches != nil && len(matches) > 0 {
			return matches[0][1]
		}
		matches = regexp.MustCompile(`^https://github.com/(.*)\n$`).FindAllStringSubmatch(string(buf), -1)
		if matches != nil && len(matches) > 0 {
			return matches[0][1]
		}
	}

	fmt.Printf("Could not determine the github repo name for this repo.\n")
	fmt.Printf("Please specify it with something like:\n")
	fmt.Printf("git config --local --add github.repo YOURUSER/YOURREPO\n")
	os.Exit(1)
	return ""
}

func CurrentUser() string {
	r, err := GithubApi("GET", "/user", nil)
	if err != nil {
		log.Fatalf("get current user: %s", err)
	}
	return r["login"].(string)
}

func GithubApiGetTyped(uri string, response interface{}) error {
	req, err := http.NewRequest("GET",
		"https://api.github.com"+uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+GithubToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 400 {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}
	return nil
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
		return nil, fmt.Errorf("%s: %s", uri, resp.Status)
	}

	body := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

func PatchLabels(issueNumber int64, addLabels []string, removeLabels []string) error {
	issueResp, err := GithubApi("GET",
		fmt.Sprintf("/repos/%s/issues/%d", GithubRepo(), issueNumber), nil)
	if err != nil {
		return err
	}

	newLabels := addLabels
	for _, l := range issueResp["labels"].([]interface{}) {
		label := l.(map[string]interface{})

		include := true
		for _, removeLabel := range removeLabels {
			if label["name"] == removeLabel {
				include = false
				break
			}
		}
		if include {
			newLabels = append(newLabels, label["name"].(string))
		}
	}
	_, err = GithubApi("PATCH",
		fmt.Sprintf("/repos/%s/issues/%d", GithubRepo(), issueNumber),
		M{"labels": newLabels})
	return err
}

package main

import (
	"fmt"
	"log"
	"os"
)

type M map[string]interface{}

func Ptal() {
	issueNumber := CurrentIssue()

	if err := Run("git", "push", "origin", CurrentBranch()); err != nil {
		os.Exit(1)
	}

	// Make sure there is a pull request for this issue (there should be)
	issueResp, err := GithubApi("GET", fmt.Sprintf("/repos/%s/issues/%d",
		GithubRepo(), issueNumber), nil)
	if err != nil {
		log.Fatalf("fetch issue: %s", err)
	}
	_, hasPullRequest := issueResp["pull_request"]
	if !hasPullRequest {
		fmt.Printf("creating pull request for issue %d\n", issueNumber)
		_, err := GithubApi("POST", fmt.Sprintf("/repos/%s/pulls", GithubRepo()),
			M{
				"issue": issueNumber,
				"head":  CurrentBranch(),
				"base":  "master",
			})
		if err != nil {
			log.Fatalf("create pull request: %s", err)
		}
	}

	// adjust lables: remove the needs-refactor, wip labels, add needs-review
	newLabels := []string{"needs-review"}
	for _, l := range issueResp["labels"].([]interface{}) {
		label := l.(map[string]interface{})
		if label["name"] == "wip" {
			continue
		}
		if label["name"] == "needs-refactor" {
			continue
		}
		newLabels = append(newLabels, label["name"].(string))
	}
	_, err = GithubApi("PATCH",
		fmt.Sprintf("/repos/%s/issues/%d", GithubRepo(), issueNumber),
		M{"labels": newLabels})
	if err != nil {
		fmt.Printf("updating issue labels: %s", err)
		os.Exit(1)
	}
}

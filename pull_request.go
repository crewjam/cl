package main

import (
	"fmt"
	"log"
	"os"
)

type PullRequestHead struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
}

type PullRequest struct {
	Head   PullRequestHead `json:"head"`
	Number int64           `json:"number"`
}

func GetPullRequestForBranch(branch string) (*PullRequest, error) {
	pullRequests := []PullRequest{}
	err := GithubApiGetTyped(
		fmt.Sprintf("/repos/%s/pulls", GithubRepo()),
		&pullRequests)
	if err != nil {
		return nil, err
	}
	for _, pullRequest := range pullRequests {
		fmt.Printf("pr.h.l=%#v (==? %#v)\n", pullRequest.Head.Label, branch)
		if pullRequest.Head.Ref == branch {
			return &pullRequest, nil
		}
	}
	return nil, nil
}

func EnsurePullRequestBranch() *PullRequest {
	currentBranch := CurrentBranch()
	pullRequest, err := GetPullRequestForBranch(currentBranch)
	if err != nil {
		log.Fatalf("get pull request: %s", err)
	}
	if pullRequest == nil {
		log.Printf("Can't find an open pull request for branch %s.\n", currentBranch)
		os.Exit(1)
	}
	return pullRequest
}

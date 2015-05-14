package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// New creates a new issue. If `arg` is numeric, it refers to an existing issue
// number. If it is a string it is the subject of a new issue to be created.
// If it is empty or not present, then we prompt you to type something.
func New(arg string) {
	if CurrentBranch() != "master" {
		fmt.Printf("The current branch is: %#v\n", CurrentBranch())
		fmt.Printf("You need to switch to master to call new.\n")
		fmt.Printf("Switch to master? (enter to continue, Ctrl+C to abort)\n")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')

		if err := Run("git", "checkout", "master"); err != nil {
			os.Exit(1)
		}
		fmt.Printf("\n")
	}

	for {
		if arg != "" {
			break
		}
		fmt.Printf("What are you about to do? Type an issue number or some ")
		fmt.Printf("text that will name a new issue\n")
		fmt.Printf("issue subject (or number): ")
		arg, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		arg = strings.TrimSpace(arg)
	}

	var issueTitle = arg
	issueNumber, err := strconv.ParseInt(arg, 10, 32)
	if err == nil {
		// Get the subject from the issue
		resp, err := GithubApi("GET",
			fmt.Sprintf("/repos/%s/issues/%d", GithubRepo(), issueNumber),
			nil)
		if err != nil {
			log.Fatalf("get issue: %s", err)
		}
		issueTitle = resp["title"].(string)
	} else {
		issueNumber = 0
	}

	branchName := strings.TrimSpace(strings.ToLower(issueTitle))
	branchName = regexp.MustCompile(`\s+`).ReplaceAllString(branchName, "-")
	branchName = regexp.MustCompile(`[^A-Za-z0-9\-]`).ReplaceAllString(branchName, "")
	if len(branchName) < 3 {
		fmt.Printf("branch name %#v seems too short\n", branchName)
		os.Exit(1)
	}
	fmt.Printf("branch: %s", branchName)

	// Create and switch to the branch
	if err := Run("git", "checkout", "-b", branchName); err != nil {
		os.Exit(1)
	}

	// Create an empty commit
	commitMessage := fmt.Sprintf("%s\n", issueTitle)
	if issueNumber > 0 {
		commitMessage += fmt.Sprintf("\nFixes #%d\n", issueNumber)
	}
	if err := Run("git", "commit", "--allow-empty", "-m", commitMessage); err != nil {
		os.Exit(1)
	}

	if err := Run("git", "push", "-u", "origin", branchName); err != nil {
		os.Exit(1)
	}

	// Create a pull request
	if issueNumber > 0 {
		fmt.Printf("creating pull request for issue %d\n", issueNumber)
		_, err := GithubApi("POST", fmt.Sprintf("/repos/%s/pulls", GithubRepo()),
			M{
				"issue": issueNumber,
				"head":  branchName,
				"base":  "master",
			})
		if err != nil {
			log.Fatalf("create pull request: %s", err)
		}
	} else {
		fmt.Printf("creating new pull request\n")
		pullReq, err := GithubApi("POST", fmt.Sprintf("/repos/%s/pulls", GithubRepo()),
			M{
				"title": issueTitle,
				"head":  branchName,
				"base":  "master",
			})
		if err != nil {
			log.Fatalf("create pull request: %s", err)
		}
		issueNumber = int64(pullReq["number"].(float64))
	}

	// adjust labels: add wip, remove needs-{review,refactor}
	if err := PatchLabels(issueNumber, []string{"wip"},
		[]string{"needs-review", "needs-refactor"}); err != nil {
		fmt.Printf("updating issue labels: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Next Steps:\n" +
		" - hack!\n" +
		" - git add / git commit / git push\n" +
		" - cl ptal\n")
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// New creates a new issue. If `arg` is numeric, it refers to an existing issue
// number. If it is a string it is the subject of a new issue to be created.
// If it is empty or not present, then we prompt you to type something.
func New(arg string) {
	if CurrentBranch() != "master" {
		fmt.Printf("The current branch is: %#v", CurrentBranch())
		fmt.Printf("You probably want to be in master to call new.\n")
		fmt.Printf("Switch to master? (enter to continue, Ctrl+C to abort)\n")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')

		if err := Run("git", "checkout", "master"); err != nil {
			os.Exit(1)
		}
	}

	issueNumber, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		issueSubject := arg
		for issueSubject == "" {
			fmt.Printf("Whacha tryin to do? I'll create an issue for you.\n")
			fmt.Printf("issue subject: ")
			issueSubject, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			issueSubject = strings.TrimSpace(issueSubject)
		}

		// create a new issue with the given subject
		resp, err := GithubApi("POST",
			fmt.Sprintf("/repos/%s/issues", GithubRepo()),
			map[string]interface{}{
				"title":    issueSubject,
				"assignee": CurrentUser(),
				"labels":   []string{"wip"},
			})
		if err != nil {
			log.Fatalf("create issue: %s", err)
		}
		issueNumber = int64(resp["number"].(float64))
		fmt.Printf("created issue #%d\n", issueNumber)
	}

	branchName := fmt.Sprintf("pr-%d", issueNumber)
	if err := Run("git", "checkout", "-b", branchName); err != nil {
		os.Exit(1)
	}

	fmt.Printf("Next Steps:\n" +
		" - hack!\n" +
		" - git add / git commit / git push\n" +
		" - cl ptal\n")
}

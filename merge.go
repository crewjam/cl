package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RevParse(rev string) string {
	buf, err := exec.Command("git", "rev-parse", rev).CombinedOutput()
	if err != nil {
		log.Fatalf("git rev-parse %s: %s", rev, err)
	}
	return strings.TrimSpace(string(buf))
}

func MergeBase(a, b string) string {
	buf, err := exec.Command("git", "merge-base", a, b).CombinedOutput()
	if err != nil {
		log.Fatalf("git merge-base %s %s: %s", a, b, err)
	}
	return strings.TrimSpace(string(buf))
}

func Merge() {
	issueNumber := CurrentIssue()
	fmt.Printf("merging issue %d to master", issueNumber)

	var commitProblems = []string{}

	// Check that the current working directory is clean
	if err := CheckWorkingDirectoryClean(); err != nil {
		commitProblems = append(commitProblems,
			fmt.Sprintf("working directory is not clean: %s", err))
	}

	// Check that the current SHA matches the remote SHA
	currentSHA := RevParse("HEAD")
	branchInfo, err := GithubApi("GET", fmt.Sprintf("/repos/%s/branches/%s", GithubRepo(), CurrentBranch()), nil)
	if err != nil {
		log.Fatalf("get branch: %s", err)
	}
	remoteCurrentSHA := branchInfo["commit"].(map[string]interface{})["sha"].(string)
	if currentSHA != remoteCurrentSHA {
		commitProblems = append(commitProblems,
			fmt.Sprintf("local HEAD is %s, remote head is %s. Do you need to push or pull changes?",
				currentSHA[:12], remoteCurrentSHA[:12]))
	}

	// Check for a +1
	commentsInfo, err := GithubApi("GET", fmt.Sprintf("/repos/%s/issues/%d/comments",
			GithubRepo(), issueNumber), nil)
	if err != nil {
		log.Fatalf("get comments: %s", err)
	}
	fmt.Printf("commentsInfo: %q\n", commentsInfo)
	hasApproval := false
	/*
	for _, ci := range commentsInfo.([]interface{}) {
		commentInfo := ci.(map[string]interface{})
		body := commentInfo["body"].(string)
		if strings.Contains(body, ":+1:") {
			hasApproval = true
		}
		if strings.Contains(body, "lgtm") {
			hasApproval = true
		}
	}
	*/
	if !hasApproval {
		commitProblems = append(commitProblems,
			fmt.Sprintf("couldn't find a :+1: or an lgtm in any comment"))
	}

	// Check that we have a passed the tests
	statusInfo, err := GithubApi("GET", fmt.Sprintf("/repos/%s/commits/%s/status",
			GithubRepo(), currentSHA), nil)
	if err != nil {
		log.Fatalf("get pull request state: %s", err)
	}
	if statusInfo["state"].(string) != "success" || statusInfo["total_count"].(float64) < 1.0 {
		commitProblems = append(commitProblems,
			fmt.Sprintf("integration status is %s (%d builders)",
				statusInfo["state"].(string),
				int(statusInfo["total_count"].(float64)) ))
	}

	// Check that the commit can be merged fast-forward without a problem.
	// Is is true iff the merge-base is also the current master
	if RevParse("master") != MergeBase("HEAD", "master") {
		commitProblems = append(commitProblems,
			fmt.Sprintf("this branch meets master at %s, but the latest master is %s. Rebase needed?",
				MergeBase("HEAD", "master")[:12], RevParse("master")[:12]))
	}

	if len(commitProblems) > 0 {
		fmt.Printf("The following problems prevent the merge:\n")
		for _, cp := range commitProblems {
			fmt.Printf(" - %s\n", cp)
		}
		os.Exit(1)
	}

	if err := Run("git", "checkout", "master"); err != nil {
		os.Exit(1)
	}
	if err := Run("git", "merge", "--ff-only", currentSHA); err != nil {
		os.Exit(1)
	}
	if err := Run("git", "push", "origin", "master"); err != nil {
		os.Exit(1)
	}

	// TODO(ross): close pull request and remove labels
	// TODO(ross): remove branch remotely
	// TODO(ross): remove branch locally
}


func CheckWorkingDirectoryClean() error {
	if err := exec.Command("git", "diff", "--exit-code").Run(); err != nil {
		return fmt.Errorf("There are unstaged changes")
	}
	if err := exec.Command("git", "diff", "--cached", "--exit-code").Run(); err != nil {
		return fmt.Errorf("There are staged, uncommitted changes")
	}
	buf, err := exec.Command("git", "ls-files", "--other", "--exclude-standard", "--directory").CombinedOutput()
	if err != nil {
		return err
	}
	if len(buf) > 0 {
		return fmt.Errorf("There are untracked files")
	}
	return nil
}

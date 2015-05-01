package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Merge() {
	branchName := CurrentBranch()
	pullRequest, err := GetPullRequestForBranch(branchName)
	if err != nil {
		log.Fatalf("get pull request: %s", err)
	}
	if pullRequest == nil {
		fmt.Printf("cannot find a pull request for branch %s\n", branchName)
		os.Exit(1)
	}

	fmt.Printf("merging PR %d to master\n.", pullRequest.Number)

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
	hasApproval, err := HasApproval(pullRequest.Number)
	if err != nil {
		log.Fatalf("get issue comments: %s", err)
	}
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
				int(statusInfo["total_count"].(float64))))
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

	// Remove the wip, needs-refactor and needs-review labels
	if err := PatchLabels(pullRequest.Number, []string{},
		[]string{"wip", "needs-refactor", "needs-review"}); err != nil {
		fmt.Printf("updating issue labels: %s", err)
	}

	// Delete the branch remotely
	_, _ = GithubApi("DELETE", fmt.Sprintf("/repos/%s/git/refs/heads/%s",
		GithubRepo(), branchName), nil)
	Run("git", "branch", "-D", branchName)
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

type GithubIssueComment struct {
	Body string `json:"body"`
}

func HasApproval(issueNumber int64) (bool, error) {
	comments := []GithubIssueComment{}
	err := GithubApiGetTyped(
		fmt.Sprintf("/repos/%s/issues/%d/comments", GithubRepo(), issueNumber),
		&comments)
	if err != nil {
		return false, err
	}
	fmt.Printf("comment: %#v\n", comments)
	for _, comment := range comments {
		if strings.Contains(comment.Body, ":+1:") {
			return true, nil
		}
		if strings.Contains(comment.Body, "lgtm") {
			return true, nil
		}
		if strings.Contains(comment.Body, "LGTM") {
			return true, nil
		}
	}
	return false, nil
}

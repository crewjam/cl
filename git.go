package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func CurrentBranch() string {
	branch, err := exec.Command("git", "symbolic-ref", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatalf("git symbolic-ref HEAD: %s", err)
	}
	return strings.TrimSpace(strings.TrimPrefix(string(branch), "refs/heads/"))
}

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

func Run(name string, args ...string) error {
	fmt.Printf("+ %s %s\n", name, strings.Join(args, " "))
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func CurrentIssue() int64 {
	branch := CurrentBranch()

	if !regexp.MustCompile(`^pr-\d+$`).MatchString(branch) {
		fmt.Printf("Oops! It looks like you are not in a PR branch\n")
		fmt.Printf("I cowardly refuse to work in a branch I don't understand\n")
		fmt.Printf("Current branch: %s", branch)
		os.Exit(1)
	}

	issueNumber, _ := strconv.ParseInt(strings.TrimPrefix(branch, "pr-"), 10, 32)
	return issueNumber
}

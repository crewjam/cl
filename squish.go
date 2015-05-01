package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Squish() {
	_ = EnsurePullRequestBranch()

	fmt.Printf("The following changes will be squished into a single commit:\n")
	cmd := exec.Command("git", "--no-pager", "log", "--format=format:%h %s", "master..HEAD")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	fmt.Printf("\nContinue?")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')

	// the commit message is the concatenation of all commits in this branch
	commitMessage, err := exec.Command("git", "log", "--reverse", "--format=format:%B", "master..HEAD").CombinedOutput()
	if err != nil {
		log.Fatalf("git log: %s", err)
	}

	Run("git", "reset", "--soft", "master")
	Run("git", "commit", "--edit", "-m", string(commitMessage))

	fmt.Printf("To push these changes:\n")
	fmt.Printf("git push origin %s --force\n", CurrentBranch())
}

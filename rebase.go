package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func MasterRevision() string {
	masterRev, err := exec.Command("git", "rev-list", "master", "-n1").CombinedOutput()
	if err != nil {
		log.Fatalf("git rev-list master: %s", err)
	}
	return strings.TrimSpace(string(masterRev))
}

func Rebase() {
	_ = CurrentIssue()

	masterRevision := MasterRevision()
	if err := Run("git", "rebase", "-i", masterRevision); err != nil {
		os.Exit(1)
	}

	fmt.Printf("To push these changes:\n")
	fmt.Printf("git push origin %s --force\n", CurrentBranch())
}

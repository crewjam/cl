package main

import (
	"fmt"
	"os"
)

func Rebase() {
	_ = CurrentIssue()

	masterRevision := RevParse("master")
	if err := Run("git", "rebase", "-i", masterRevision); err != nil {
		os.Exit(1)
	}

	fmt.Printf("To push these changes:\n")
	fmt.Printf("git push origin %s --force\n", CurrentBranch())
}

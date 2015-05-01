package main

import (
	"fmt"
	"os"
)

type M map[string]interface{}

func Ptal() {
	currentBranch := CurrentBranch()
	pullRequest := EnsurePullRequestBranch()

	if err := Run("git", "push", "origin", currentBranch); err != nil {
		os.Exit(1)
	}

	// adjust labels: remove the needs-refactor, wip labels, add needs-review
	if err := PatchLabels(pullRequest.Number, []string{"needs-review"},
		[]string{"wip", "needs-refactor"}); err != nil {
		fmt.Printf("updating issue labels: %s", err)
		os.Exit(1)
	}
}

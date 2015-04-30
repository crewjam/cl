package main

// Save stages all changes, commits them, and pushes them with a generic commit
// message. Use it to keep from losing stuff.
func Save() {
	_ = CurrentIssue()

	Run("git", "add", "-A")
	Run("git", "commit", "-a", "-m", "automatic progress commit")
	Run("git", "push", "origin", CurrentBranch())
}

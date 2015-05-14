package main

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Printf(`usage: %s <subcommand> <options>

Subcommands:

  help   - display this text
  new    - create a new change and switch to it
  squish - combine commits
  rebase - combine commits using rebase
  save   - commit all changes with a generic commit message
  ptal   - mark a PR as ready for review
  merge  - merge a PR into the master

`, os.Args[0])
}

func main() {
	operation := ""
	if len(os.Args) > 1 {
		operation = os.Args[1]
	}

	switch operation {
	case "new":
		arg := ""
		if len(os.Args) > 2 {
			arg = os.Args[2]
		}
		New(arg)
		os.Exit(0)
	case "squish":
		Squish()
		os.Exit(0)
	case "rebase":
		Rebase()
		os.Exit(0)
	case "save":
		Save()
		os.Exit(0)
	case "pull":
	case "merge":
		Merge(os.Args[2:])
		os.Exit(0)
	case "ptal":
		Ptal()
	case "help":
		Usage()
		os.Exit(0)
	default:
		fmt.Printf("unknown operation: %s", operation)
		Usage()
		os.Exit(1)
	}
}

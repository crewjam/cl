#!/bin/bash







# save
#
# stages all changes, commits them, and pushes them with a generic commit
# message.
save() {
    if [[ -z $(current_pr) ]] ; then
      echo "Oops! I'm lame and I won't work with non-pr branches"
      echo "current branch: $(current_branch)"
      echo "current pr: $(current_pr)"
      exit 1
    fi

    (
        set -ex
        git add -A
        git commit -a -m "automatic progress commit"
        git push origin $(current_branch)
    )
}


check_working_directory_clean() {
    if ! git diff --exit-code > /dev/null ; then
        git status
        echo "cl: working directory is not clean"
        exit 1
    fi

    if ! git diff --cached --exit-code > /dev/null ; then
        git status
        echo "cl: working directory is not clean"
        exit 1
    fi

    if ! [[ -z $(git ls-files --other --exclude-standard --directory) ]] ; then
        git status
        echo "cl: working directory is not clean"
        exit 1
    fi
}

merge() {
    pr=$(current_pr)
    if [[ -z $pr ]] ; then
      echo "Oops! I'm lame and I won't work with non-pr branches"
      echo "current branch: $(current_branch)"
      echo "current pr: $pr"
      exit 1
    fi

    check_working_directory_clean || exit 1

    # look for a +1
    comments=$(curl -H "Authorization: token $github_token" \
      "https://api.github.com/repos/$github_repo/issues/$pr/comments" |\
       jq -r '.[].body')
    if ! echo $comments | grep ":+1:" >/dev/null ; then
        if ! echo $comments | grep "lgtm" >/dev/null ; then
            echo "pull request is missing :+1: or lgtm"
            exit 1
        fi
    fi

    # check that the stuff works
    mergeable_state=$(curl -H "Authorization: token $github_token" \
      "https://api.github.com/repos/$github_repo/pulls/$pr" |\
      jq -r '.mergeable_state')
    if [[ "$mergeable_state" != "clean" ]] ; then
      echo "pr is ${mergeable_state}. Did the tests pass?"
      exit 1
    fi

    # TODO(ross): check that our current SHA is the correct SHA on the server,
    # bail if not

    git checkout master
    git merge --ff-only

    echo "not yet implemented"
    exit 1
}

main() {
    #set -ex
    operation=$1
    shift
    case "$operation" in
        new)
            new $*
            ;;
        squish)
            squish $*
            ;;
        rebase)
            rebase $*
            ;;
        save)
            save $*
            ;;
        ptal)
            ptal $*
            ;;
        merge)
            merge $*
            ;;
        help)
            usage
            exit 0
            ;;
        *)
            echo "unknown operation: $operation"
            echo ""
            usage
            exit 1
            ;;
    esac
}

main $*
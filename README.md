# cl is a tool to make github pull requests act more like gerrit.


## Create a new CL:

    $ cl new
    7q8

What this does:

    $ git branch 7q8
    $ git checkout 7q8

## Create a PR:

    $ git add ...
    $ git commit ...
    $ cl push

What this does:

    $ git push 7q8 
    $ curl ... api.github.com/repos/xxx/yyy/pulls -d '{title: ..., head: 7q8, base: master}'

## Review state

Review state is encoded in the topic of the pull request.

    # Review state (please merge with `cl merge`)
    R: ross +1
    R: jim
    R: cathy -1
    V: lint +1
    V: test +1

## Mark a review:

   cl 7q8 -1

## Merge a PR:

    cl merge

What this does:

  - in a fresh repo
  - check that the merge criteria are met
  - squish the commits into one, rewriting the commit message
  - merge


## Show state:

    cl open


 



A random CL name is [0-9a-z]{3} (which is > 46k changes)

Things to check:

 - CLs that depend on other CLs

    


# cl is a tool for working with GitHub pull requests

This tool assumes you have a workflow where commits generally enter the master
branch via pull requests. 

Note: "cl" == "change list" which is the term 
[Gerrit](http://lmgtfy.com/?q=gerrit+code+review) uses to describe a related
set of changes.

To start working on an existing issue:

    $ cl new 42
    
To start working creating a new issue / pull request:

    $ cl new "my fancy new feature"
    
This will create a new issue to track the change, a branch with a single empty
commit in it, and a pull request.
 
You can commit away as normal, pushing changes as you go. Then when you are 
ready, you can create the pull request and ask for review:
 
    $ cl ptal
     
This removes the `wip` label and adds the `review-needed` label.

You may wish to merge your changes into fewer commits. To do this you can either
squish or rebase. 

    $ cl squish
    
This command resets ``.git`` so all the changes in your branch look like they 
are uncommitted against the master. Then it invokes commit and you can commit
them as one commit. Finally you'll need to push to the origin with the 
``--force`` flag (because you are "rewriting history").

    $ cl rebase
    
This works similarly. It invokes ``git rebase -i`` between the master
and your branch.

When all is ready to go, you are ready to merge:

    $ cl merge
    
This checks that:

  - your code can be fast-forwarded (that you don't need to rebase your branch)
  - that you've got a ``:+1:`` or ``lgtm`` in your PR comment.
  - that there is at least one CI service commenting on your commit and that it 
    is a success.

If all this is true, then it merges your branch, pushes the code and closes the
branch and pull request.

Happy Hacking!
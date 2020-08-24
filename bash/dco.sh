#!/bin/bash

function usage() {
    echo "ERROR: -h option provided or given path does not exist"
    echo "Usage:"
    echo "  ./$(basename $0) [-h|<repository path with exists>]"
    echo "E.g."
    echo "  ./$(basename $0) /opt/github/user1/tmp_repo"
    exit -1
}

function main {
    # Usage
    if [ "$1" = "-h" ] || [ ! -d "$1"  ]; then
        usage;
    fi

    # Perform operation
    (
        # Check correct author
        cd "$1";
        user=$(git config --get user.name)
        echo "Going to add DCO to all commits by user: '$user'. Okay? [y/n]"
        read prompt
        if [ "$prompt" != "y" ] && [ "$prompt" != "Y" ]; then
            echo "INFO: n selected, Aborting..."
            exit 0
        fi

        # Apply changes
        echo "INFO: Applying DCO for user '$user'"
        first_commit=$(git rev-list --max-parents=0 HEAD | tail -n 1)
        for i in $(git log --oneline --author="$user"|awk {'print $1'}); do
            GIT_SEQUENCE_EDITOR="sed -i -re \"s/^pick $i/e $i/\"" git rebase -i $first_commit 
            git commit --amend -s --no-edit -q
            echo "INFO: Added sign-off for commit: $i"
            git rebase --continue
        done
        echo "INFO: Performing git push..."
        git push origin master -f
    )

}

main "$@"
#!/usr/bin/env bash

set -e

if [[ ! -f "build/env.sh" ]]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="${PWD}/build/_workspace"
root="${PWD}"
berithDir="${workspace}/src/github.com/BerithFoundation"
projectName="berith-chain"
if [[ ! -L "${berithDir}/${projectName}" ]]; then
    mkdir -p "${berithDir}"
    cd "${berithDir}"
    ln -s ../../../../../. ${projectName}
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "${berithDir}/${projectName}"
PWD="${berithDir}/${projectName}"

# Launch the arguments with the configured environment.
exec "$@"

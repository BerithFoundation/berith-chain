#!/bin/sh

set -e

if [ ! -f "build.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
berdir="$workspace/src/github.com/ibizsoftware/berith-chain"
if [ ! -L "$berdir/berith" ]; then
    mkdir -p "$berdir"
    cd "$berdir"
    ln -s ../../../../../. go-ethereum
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$berdir/berith"
PWD="$berdir/berith"

# Launch the arguments with the configured environment.
exec "$@"
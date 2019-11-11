#!/bin/sh

set -e

workspace="${PWD}/_workspace"
root="${PWD}"
berithDir="${workspace}/src/github.com/BerithFoundation"
projectName="berith-chain"
gobin="${PWD}/_bin"
bindir="${PWD}/bin"

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

GOBIN="$gobin"
export GOBIN

echo "#################################################################"
echo "#######              Set up go-astilectron-bundler     ##########"
echo "#################################################################"

go get -u github.com/asticode/go-astilectron-bundler/...
go install github.com/asticode/go-astilectron-bundler/astilectron-bundler


echo "#################################################################"
echo "#######              Set up berith-chain               ##########"
echo "#################################################################"

if [ ! -L "${berithDir}/${projectName}" ]; then
    mkdir -p "${berithDir}"
    cd "${berithDir}"
    ln -s ../../../../../. ${projectName}
    cd "$root"
fi

echo "#################################################################"
echo "#######              Build berith wallet               ##########"
echo "#################################################################"

cd ${berithDir}/${projectName}/wallet && astilectron-bundler -o $bindir
rm -rf $gobin
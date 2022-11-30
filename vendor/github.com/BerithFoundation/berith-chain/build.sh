#!/usr/bin/env bash

# absolute current path
SCRIPT_PATH=$(cd "$(dirname $0)" && pwd)
MODULE=
MODULES=

function printHelp() {
  echo "This command is for building berith"
  echo "Below is available commands and options."
  echo ""
  echo "Usage"
  echo "build.sh [MODULE]"
  echo ""
  echo "Modules :"
  echo "  berith                     build a berith client"
  echo "Options :"
  echo "  -h, --help                 show available commands and options."
  echo ""
}

function buildBerith() {
    echo "#################################################################"
    echo "#######              Build berith module               ##########"
    echo "#################################################################"
    rm -rf ${SCRIPT_PATH}/build/bin/berith
    go build -o ${SCRIPT_PATH}/build/bin/berith ${SCRIPT_PATH}/cmd/berith/
    res=$?
    set +x
    if [ ${res} -ne 0 ]; then
        echo "Failed to build berith module"
        exit 1
    fi
}


# parse module
if [[ -z ${1} ]]; then
    MODULE="all"
else
    case "${1}" in
      berith | all)
        MODULE=${1}
        ;;
      -h | --help )
        printHelp
        exit 0
        ;;
      * )
        echo "Invalid module name : ${1}"
        echo ""
        printHelp
        exit 1
        ;;
    esac
    shift
fi

if [[ ${MODULE} == "all" ]]; then
    MODULES=("berith")
else
    MODULES=(${MODULE})
fi

for M in ${MODULES[@]}; do
    case "${M}" in
      berith )
        buildBerith
        ;;
    esac
done


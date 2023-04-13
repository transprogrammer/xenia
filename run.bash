#!/usr/bin/env bash

# REQ: Runs xenia. <>

set -o errexit
set -o xtrace

cd src
go run .

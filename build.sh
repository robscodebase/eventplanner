#!/bin/bash
# Exit shell on error.
set -e
#cd src/server
#go test
#cd ../../
mkdir -p bin && cd bin/
go build -v ../src/server



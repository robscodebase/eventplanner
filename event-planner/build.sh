#!/bin/bash
# Exit shell on error.
set -e

mkdir -p bin && cd bin/
go build -v ../src/server



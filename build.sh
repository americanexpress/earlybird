#!/bin/sh

#
# Copyright 2023 American Express
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing
# permissions and limitations under the License.
#

# Change to directory of build script to support relative paths
cd $(dirname "$0")

# Build Version
if [ -z "$1" ]
then
      EBVERSION='dev'
else
      EBVERSION="$1"
fi

# Start
echo
echo "Building EarlyBird - version ${EBVERSION}"
echo "-----------------------------------------------------------------"
echo

# Tidy up and download depdendencies
echo "[x] Running go mod tidy for dependency maintenance"
go mod tidy
echo "COMPLETE"
echo


# Run Unit Tests
############################
# Unit test outcomes are stored to a file that is read for falures, then deleted afterwards.

echo "[x] Running Unit Tests"

# Set log location for unit test outcomes
log_file="${HOME}/eb-build-log.log"

# Run Unit Tests
go test -p 10 ./pkg/... -covermode=count > ${log_file}

# Check Unit Test Outcomes; exit 1 if failure detected!
testOutput=$? #get the return value of the last executed command
if [ ${testOutput} -ne 0 ]; then
  echo "FAILED"
  echo
  echo "********************LOG********************"
  cat ${log_file}
  echo "********************LOG********************"
  echo
  echo "exit 1"
  rm -rf ${log_file}
  exit 1
else 
  cat ${log_file}
  rm -rf ${log_file}
  echo "PASSED"
fi
echo
############################

# If the above check went okay, test cases have passed and the script can now proceed to build binaries
echo "[x] Building Linux binary"
env GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=${EBVERSION}'" -o binaries/go-earlybird-linux
echo "COMPLETE"
echo

echo "[x] Building MacOS binary"
env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=${EBVERSION}'" -o binaries/go-earlybird-mac
echo "COMPLETE"
echo

echo "[x] Building Windows binary"
env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=${EBVERSION}'" -o binaries/go-earlybird.exe
echo "COMPLETE"
echo

# Completion
echo "-----------------------------------------------------------------"
echo "Build Finished. Binary location: $(pwd)/binaries"
echo
#!/bin/sh

#
# Copyright 2021 American Express
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

if [ -z "$1" ]
then
      echo "Earlybird version is empty. Setting to default 'dev'"
      version='dev'
else
      echo "Earlybird version being build $1"
      version="$1"
fi

echo "Running Unit Tests and Building binaries... version: $version"
echo "-----------------------------------------------------------------"
echo
echo "Tidy and Download modules ..."
go mod tidy && go mod download -x

echo "Running Unit Tests..."
go test -p 10 ./pkg/** -covermode=count
testOutput=$? #get the return value of the last executed command
#exit with non-Zero to stop the build if unit tests fail
if [ $testOutput -ne 0 ]; then
  echo "EXIT... Test Failed" >&2
  exit $testOutput
fi

echo "Building Linux binary"
env GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=$version'" -o binaries/go-earlybird-linux
echo "Building Linux binary - Completed!!!"
echo "Building Windows binary"
env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=$version'" -o binaries/go-earlybird.exe
echo "Building Windows binary - Completed!!!"
echo "Building MacOS binary"
env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=$version'" -o binaries/go-earlybird
echo "Building MacOS binary - Completed!!!"
echo "Build Completed ..."

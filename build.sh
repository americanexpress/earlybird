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

echo "Running Unit Tests"
tmpfile=/tmp/eb-tmp.log
if SCRIPT_OUTPUT=$(local=true go test ./pkg/... > ${tmpfile}); then
    cat ${tmpfile}
	
    echo "Running Go FMT"
    go fmt ./pkg/...
    
    mkdir -p binaries

    echo "Building Linux binary"
    env GOOS=linux GOARCH=amd64 go build -o binaries/go-earlybird-linux

    echo "Building Windows binary"
    env GOOS=windows GOARCH=amd64 go build -o binaries/go-earlybird.exe

    echo "building MacOS binary"
    env GOOS=darwin GOARCH=amd64 go build -o binaries/go-earlybird

    echo "Cross-Compilation Complete"
else
    echo "Unit Tests FAILED!"
    cat ${tmpfile}
fi

rm $tmpfile

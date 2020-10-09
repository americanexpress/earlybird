::
:: Copyright 2020 American Express
::
:: Licensed under the Apache License, Version 2.0 (the "License");
:: you may not use this file except in compliance with the License.
:: You may obtain a copy of the License at
::
:: http://www.apache.org/licenses/LICENSE-2.0
::
:: Unless required by applicable law or agreed to in writing, software
:: distributed under the License is distributed on an "AS IS" BASIS,
:: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
:: or implied. See the License for the specific language governing
:: permissions and limitations under the License.
::
@ECHO off
ECHO Running Unit Tests
SET tmpfile=c:\temp\eb-tmp.log
TYPE NUL > blank
go test ./pkg/... > %tmpfile% -mod vendor
FC tmpfile blank >NUL
if errorlevel 0 (
    TYPE %tmpfile%
	
    ECHO Running Go FMT
    go fmt ./pkg/...

    ECHO Building Windows binary
    SET GOOS=windows
    SET GOARCH=amd64
    go build -o binaries/go-earlybird.exe -mod vendor

    ECHO Cross-Compilation Complete
) else (
    ECHO Unit Tests FAILED!
    TYPE %tmpfile%
)

DEL %tmpfile%

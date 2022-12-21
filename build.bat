@ECHO off

::
:: Copyright 2023 American Express
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

:: Change to directory of build script to support relative paths
CD /D %~dp0

:: Build Version
if "%1" == "" (
    SET EBVERSION=dev
) else (
    SET EBVERSION=%1
)

:: Start
ECHO.
ECHO Building EarlyBird
ECHO -----------------------------------------------------------------
ECHO.

:: Tidy up and download depdendencies
ECHO [x] Running go mod tidy for dependency maintenance
go mod tidy
ECHO COMPLETE
ECHO.

:: Run Unit Tests
::::::::::::::::::::::::::::
:: Unit test outcomes are stored to a file that is read for falures, then deleted afterwards.

ECHO [x] Running Unit Tests

:: Set log location for unit test outcomes
SET log_file=%APPDATA%\eb-build-log.log

:: Delete log to start fresh
IF EXIST %log_file% DEL %log_file%

:: Run Unit Tests
go test -p 10 .\pkg\... -covermode=count > %log_file% 

:: Check Unit Test Outcomes; exit 1 if failure detected!
FINDSTR /m "FAIL" %log_file% > NUL
if %ERRORLEVEL% == 0 (
    ECHO FAILED
    ECHO.
    ECHO ********************LOG********************
    TYPE %log_file%
    ECHO ********************LOG********************
    ECHO.
    ECHO exit 1
    ECHO.
    DEL %log_file%
    exit(1)
) else (
    TYPE %log_file%
    DEL %log_file%
    ECHO PASSED
)
ECHO.
::::::::::::::::::::::::::::

:: If the above check went okay, test cases have passed and the script can now proceed to build binaries
ECHO [x] Building Linux binary
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=%EBVERSION%'" -o binaries/go-earlybird-linux
ECHO COMPLETE
ECHO.

ECHO [x] Building MacOS binary
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=%EBVERSION%'" -o binaries/go-earlybird-mac
ECHO COMPLETE
ECHO.

ECHO [x] Building Windows binary
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=%EBVERSION%'" -o binaries/go-earlybird.exe
ECHO COMPLETE
ECHO.

:: Completion
ECHO -----------------------------------------------------------------
ECHO Build Finished. Binary location: %~dp0binaries\
ECHO.

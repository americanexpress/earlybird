::
:: Copyright 2021 American Express
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
ECHO Installing Go-EarlyBird

SET ebPath=%USERPROFILE%\AppData\go-earlybird\
IF NOT EXIST %ebPath% MKDIR %ebPath%

ECHO Copying executable
COPY binaries\go-earlybird.exe %ebPath% >NUL
ECHO Copying configurations
COPY config\* %ebPath% >NUL
COPY .ge_ignore %ebPath% >NUL

REM Make sure AppData\go-earlybird is in the path
ECHO Updating path
SET "PATH=%PATH%;%ebPath%"

ECHO Setup completed.  Start a new CMD session for the path change to take effect.

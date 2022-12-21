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

:: Start
ECHO.
ECHO Uninstalling EarlyBird
ECHO -----------------------------------------------------------------
ECHO.

:: Installation paths
SET ebBuildLog=%APPDATA%\eb-build-log.log
SET ebPath=%APPDATA%\go-earlybird

:: Delete left-over build logs
ECHO [x] Removing %ebBuildLog% if it exists
IF EXIST %ebBuildLog% DEL /S %ebBuildLog%
ECHO COMPLETE
ECHO.

:: Delete config/install directory
ECHO [x] Removing %ebPath% if it exists
IF EXIST %ebPath% RD /S /Q %ebPath%
ECHO COMPLETE
ECHO.

:: Updating PATH
::::::::::::::::::::::::::::
:: Dilligence was spent to ensure SYTEM PATH and USER PATH environment variables are not merged.

ECHO [x] Ensuring no traces of %ebPath% in PATH user environment variable

:: Query registry key for current user path registry key and set its data to userPath
for /f "tokens=3" %%A in ('REG QUERY "HKEY_CURRENT_USER\Environment" /v "PATH"') DO SET userPath=%%A

:: Remove go-earlybird installation directory to PATH user environment variable
CALL SET newPath=%%userPath:%ebPath%=%%

:: Strip double semi-colons
SET newPath=%newPath:;;=;%

:: Replace PATH user environment variable with newPath variable built.
REG ADD HKEY_CURRENT_USER\Environment /f /v Path /t REG_EXPAND_SZ /d %newPath%

ECHO COMPLETE
ECHO.
::::::::::::::::::::::::::::

:: Complete
ECHO -----------------------------------------------------------------
ECHO Uninstall Finished.
ECHO.

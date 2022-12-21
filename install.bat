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

:: Change to directory of install script to support relative paths
CD /D %~dp0

:: Start
ECHO.
ECHO Installing EarlyBird
ECHO -----------------------------------------------------------------
ECHO.

:: Installation path
SET ebConfigPath=%APPDATA%\go-earlybird

:: Ensure install directory is created
ECHO [x] Creating install path (%ebConfigPath%\)
IF NOT EXIST %ebConfigPath% MKDIR %ebConfigPath%
ECHO COMPLETE
ECHO.

:: Copy configurations
ECHO [x] Copying configurations to %ebConfigPath%\
XCOPY /y /s config\* %ebConfigPath% > NUL
ECHO COMPLETE
ECHO.

:: Copy .ge_ignore
ECHO [x] Copying .ge_ignore to %ebConfigPath%\
COPY .ge_ignore %ebConfigPath% > NUL
ECHO COMPLETE
ECHO.

:: Copy executable
ECHO [x] Copying go-earlybird executable to %ebConfigPath%\
COPY binaries\go-earlybird.exe %ebConfigPath% > NUL
ECHO COMPLETE
ECHO.

:: Updating PATH
::::::::::::::::::::::::::::
:: Dilligence was spent to ensure SYTEM PATH and USER PATH environment variables are not merged.

ECHO [x] Updating PATH user environment variable with go-earlybird installation directory if not already in place

:: Query registry key for current user path registry key and set its data to userPath
for /f "tokens=3" %%A in ('REG QUERY "HKEY_CURRENT_USER\Environment" /v "PATH"') DO SET userPath=%%A

:: Set compareUserPathNoEB to userPath with any potential references to ebConfigPath stripped, removing semi-colons for future string comparisons
CALL SET compareUserPathNoEB=%%userPath:%ebConfigPath%=%%
SET compareUserPathNoEB=%compareUserPathNoEB:;=%

:: Set compareUserPath to userPath to be later compared with compareUserPathNoEB, removing semi-colons for future string comparisons
SET compareUserPath=%userPath:;=%

:: Check if userPath contains the go-earlybird install directory; if not; add it!
if %compareUserPathNoEB%==%compareUserPath% GOTO set_user_path
GOTO END

:set_user_path
:: Append go-earlybird installation directory to PATH user environment variable
SET newPath="%userPath%;%ebConfigPath%"
:: Strip any double semi-colons
SET newPath=%newPath:;;=;%
:: Replace PATH user environment variable with newPath variable built.
REG ADD HKEY_CURRENT_USER\Environment /f /v Path /t REG_EXPAND_SZ /d %newPath%
GOTO END
:END

ECHO COMPLETE
ECHO.
::::::::::::::::::::::::::::

:: Complete
ECHO -----------------------------------------------------------------
ECHO Install Finished. 
ECHO Config location: %ebConfigPath%\
ECHO Binary location: %ebConfigPath%\go-earlybird.exe
ECHO.

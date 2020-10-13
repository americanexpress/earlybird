#!/bin/sh

#
# Copyright 2020 American Express
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

echo Installing Go-EarlyBird
# create the default config directory structure
CONFIGDIR="$HOME/.go-earlybird"
LOCALBINDIR="/usr/local/bin"
if [ ! -d "$CONFIGDIR" ]; then
  mkdir -p "$CONFIGDIR"
fi

# Pull the binary
echo Copying executable
if [[ "$OSTYPE" == "darwin"* ]]; then
        #Copy Mac OS bin
        cp binaries/go-earlybird $LOCALBINDIR/go-earlybird
        else
        #Copy Linux bin
        cp binaries/go-earlybird-linux $LOCALBINDIR/go-earlybird
        fi

echo Updating permissions
chmod u+x $LOCALBINDIR/go-earlybird

# Pull the module configs
echo Copying configurations
cp -a config/. $CONFIGDIR
cp .ge_ignore $HOME/.ge_ignore
echo Setup complete

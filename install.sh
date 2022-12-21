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

# Change to directory of install script to support relative paths
cd $(dirname "$0")

# Start
echo
echo "Installing EarlyBird"
echo "-----------------------------------------------------------------"
echo

# Installation paths
ebConfigPath="${HOME}/.go-earlybird/"
localBinDir="/usr/local/bin"

# Ensure install directory is created
echo "[x] Creating install path (${ebConfigPath})"
if [ ! -d "${ebConfigPath}" ]; then
  mkdir -p "${ebConfigPath}"
fi
echo "COMPLETE"
echo

# Copy configurations
echo "[x] Copying configurations to ${ebConfigPath}"
cp -r config/* ${ebConfigPath}
echo "COMPELTE"
echo

# Copy .ge_ignore
echo "[x] Copying .ge_ignore to ${ebConfigPath}"
cp .ge_ignore ${ebConfigPath}
echo "COMPLETE"
echo

# Copy executable
echo "[x] Copying go-earlybird executable to ${localBinDir}/"
if [ "${OSTYPE}" = "darwin" ]; then
  #Copy Mac OS bin
  cp binaries/go-earlybird-mac ${localBinDir}/go-earlybird
else
  #Copy Linux bin
  cp binaries/go-earlybird-linux ${localBinDir}/go-earlybird
fi
echo "COMPLETE"
echo

# Updating binary permissions
echo "[x] Updating permissions to ${localBinDir}/go-earlybird"
chmod u+x ${localBinDir}/go-earlybird
echo "COMPLETE"
echo

# Complete
echo "-----------------------------------------------------------------"
echo "Install Finished."
echo "Config location: ${ebConfigPath}"
echo "Binary location: ${localBinDir}/go-earlybird"
echo

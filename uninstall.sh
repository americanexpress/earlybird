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

# Start
echo
echo "Uninstalling EarlyBird"
echo "-----------------------------------------------------------------"
echo

# Installation paths
ebBuildLog="${HOME}/eb-build-log.log"
ebConfigPath="${HOME}/.go-earlybird"
localBinDir="/usr/local/bin"

# Delete left-over build logs
echo "[x] Removing ${ebBuildLog} if it exists"
rm -rf ${ebBuildLog}
echo "COMPLETE"
echo

# Delete config directory
echo "[x] Removing ${ebConfigPath} if it exists"
rm -rf ${ebConfigPath}
echo "COMPLETE"
echo

# Delete binary
echo "[x] Removing ${localBinDir}/go-earlybird if it exists"
rm -rf "${localBinDir}/go-earlybird"
echo "COMPLETE"
echo

# Complete
echo "-----------------------------------------------------------------"
echo "Uninstall Finished."
echo

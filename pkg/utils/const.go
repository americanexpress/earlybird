/*
 * Copyright 2023 American Express
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package utils

const (
	errInvalidPath   string = "Invalid Path. Exiting."
	ebConfFileDir    string = ".go-earlybird"
	ebConfFileName   string = "earlybird.json"
	ebWinConfFileDir string = "\\AppData\\Roaming\\go-earlybird"
	gitHTTP          string = "http://"
	gitHTTPS         string = "https://"
	gitPasswdPrompt  string = "Enter your git password: "
	errGitPasswd     string = "Failed to receive git password input"
	errGitDelete     string = "Failed to delete git dir: "
	//Staged is the const for tracking staged files
	Staged string = "staged"
	//Tracked is the const for tracking tracked files
	Tracked string = "tracked"
	//All is the const for tracking all files
	All string = "all"
)

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

package main

import (
	"flag"
	"os"

	"github.com/americanexpress/earlybird/pkg/core"
	"github.com/americanexpress/earlybird/pkg/scan"
	"github.com/americanexpress/earlybird/pkg/utils"
)

var (
	eb     core.EarlybirdCfg
	ptr    core.PTRHTTPConfig
	gitcfg core.PTRGitConfig
)

func main() {
	//Define HTTP server cli params
	ptr.HTTP = flag.String("http", "", "Listen IP and Port for HTTP API e.g. 127.0.0.1:8080")
	ptr.HTTPConfig = flag.String("http-config", "", "Path to webserver config JSON file")
	ptr.HTTPS = flag.String("https", "", "Listen IP and Port for HTTPS/2 API e.g. 127.0.0.1:8080 (Don't forget the https-cert and https-key flags)")
	ptr.HTTPSCert = flag.String("https-cert", "", "Certificate file for TLS")
	ptr.HTTPSKey = flag.String("https-key", "", "Private key file for TLS")
	//Define Git cli params
	gitcfg.Project = flag.String("git-project", "", "Full URL to a github organization to scan e.g. github.com/org")
	gitcfg.Repo = flag.String("git", "", "Full URL to a git repo to scan e.g. github.com/user/repo")
	gitcfg.RepoUser = flag.String("git-user", os.Getenv("gituser"), "If the git repository is private, enter an authorized username")
	gitcfg.RepoBranch = flag.String("git-branch", "", "Name of branch to be scanned")

	//Load CLI params and Earlybird config
	eb.ConfigInit()

	//Load in the rules
	scan.Init(eb.Config)

	//Start earlybird webserver when parameters provided
	if *ptr.HTTP != "" || *ptr.HTTPS != "" {
		eb.StartHTTP(ptr)
	}

	//Clone the repository we're scanning when param provided
	eb.GitClone(gitcfg)

	//Scan delete anything
	eb.Scan()

	//Delete any left over git data
	utils.DeleteGit(*gitcfg.Repo, eb.Config.SearchDir)
}

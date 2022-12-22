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

package git

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/americanexpress/earlybird/pkg/utils"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// ReposPerProject returns all the repositories contained within a  bitbucket or github project
func ReposPerProject(projectURL, username, password string) (scanRepos []string) {
	if strings.Contains(projectURL, "github.com/") { //Scan Github
		var basicauth github.BasicAuthTransport
		basicauth.Username = username
		basicauth.Password = password

		client := github.NewClient(basicauth.Client())
		// list public repositories for org "github"
		opt := &github.RepositoryListByOrgOptions{Type: "public"}
		repos, _, err := client.Repositories.ListByOrg(context.Background(), utils.GetGitProject(projectURL), opt)
		if err != nil {
			log.Println("Failed To Get Project Repositories:", err)
			os.Exit(1)
		}
		for _, repo := range repos {
			scanRepos = append(scanRepos, *repo.HTMLURL)
		}
	} else { //Scan bitbucket
		baseurl, path, project := utils.ParseBBURL(projectURL)
		client := newBitClient(baseurl, path, username, password)

		requestParams := pagedRequest{
			Limit: 100000,
			Start: 0,
		}

		reposResponse, err := client.getRepositories(project, requestParams)
		if err != nil {
			log.Println("Failed To Get Project Repositories:", err)
			os.Exit(1)
		}

		for _, repo := range reposResponse.Values {
			links := repo.Links["clone"]
			for _, link := range links {
				if link["name"] == "http" {
					scanRepos = append(scanRepos, link["href"])
				}
			}
		}
	}
	return scanRepos
}

// CloneGitRepos Clones a Git repo into a random temporary folder
func CloneGitRepos(repoURLs []string, username, password string, branch string, json bool) (tmpDir string, err error) {
	tmpDir, err = ioutil.TempDir("", "ebgit")
	if err != nil {
		return "", err
	}
	auth := &http.BasicAuth{
		Username: username,
		Password: password,
	}

	for _, repo := range repoURLs {
		options := git.CloneOptions{
			URL:   repo,
			Depth: 1,
		}

		if username != "" {
			options.Auth = auth
		}

		if branch != "" {
			options.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
		}

		if !json {
			log.Println("Cloning Repository:", repo)
			if branch != "" {
				log.Println("Cloning Branch:", branch)
			}
			options.Progress = os.Stdout
		}

		scanDir := tmpDir + "/" + utils.GetGitRepo(repo)
		if len(repoURLs) == 1 {
			scanDir = tmpDir
		}

		//Clone repo into random temporary path
		log.Println("Cloned into:", scanDir)
		_, err = git.PlainClone(scanDir, false, &options)
		if err != nil {
			return tmpDir, err
		}
	}
	return tmpDir, err
}

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
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

const baseURI = "/rest/api/1.0"

// BBCommitAPIURL formats multiple values into the bitbucket REST API URL
func BBCommitAPIURL(host, source, project, repository, commitid string) (url string) {
	//default protocol is https unless server is local
	protocol := "https://"
	if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
		protocol = "http://"
	}
	return protocol + host + "/" + source + "/rest/api/latest/projects/" + project + "/repos/" + repository + "/commits/" + commitid
}

// BBCommitURL formats multiple values into the bitbucket commit URL
func BBCommitURL(host, source, project, repository, commitid string) (url string) {
	return "https://" + host + "/" + source + "/projects/" + project + "/repos/" + repository + "/commits/" + commitid
}

func newBitClient(url string, path string, username string, password string) *bitClient {

	return &bitClient{
		sling: sling.New().Base(url).Set("Accept", "application/json").SetBasicAuth(username, password),
		Path:  path,
	}
}

func (bc *bitClient) getRepositories(projectKey string, params pagedRequest) (getRepositoriesResponse, error) {

	response := getRepositoriesResponse{}

	_, err := bc.doGet(bc.Path,
		fmt.Sprintf("/projects/%s/repos", projectKey),
		params,
		&response,
	)

	return response, err
}

func (bc *bitClient) doGet(path, uri string, params interface{}, rData interface{}) (*http.Response, error) {

	rError := new(errorResponse)

	resp, _ := bc.sling.New().Get(path+baseURI+uri).QueryStruct(params).Receive(rData, rError)

	return bc.checkReponse(resp, rError)
}

func (bc *bitClient) checkReponse(resp *http.Response, errorResponse *errorResponse) (*http.Response, error) {

	if resp != nil && resp.StatusCode > 299 {

		message := fmt.Sprintf("%s - %s\n", resp.Status, resp.Request.URL.String())
		for _, e := range errorResponse.Errors {
			message += e.Context + ": " + e.Message + "\n"
		}

		return nil, requestError{
			Code:    resp.StatusCode,
			Message: message,
		}
	}
	return resp, nil
}

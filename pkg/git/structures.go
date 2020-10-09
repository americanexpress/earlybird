/*
 * Copyright 2020 American Express
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

import "github.com/dghubble/sling"

type errorResponse struct {
	Errors []apiError
}

type pagedRequest struct {
	Limit uint `url:"limit,omitempty"`
	Start uint `url:"start,omitempty"`
}

type pagedResponse struct {
	Size       uint `json:"size"`
	Limit      uint `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Start      uint `json:"start"`
}

type link map[string]string
type links map[string][]link

type project struct {
	Key         string `json:"key,omitempty"`
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public,omitempty"`
	Type        string `json:"type,omitempty"`
	Links       links  `json:"links,omitempty"`
}

type repository struct {
	Slug          string
	ID            uint
	Name          string
	ScmID         string
	State         string
	StatusMessage string
	Forkable      bool
	Project       project
	Public        bool
	Links         links
}

type apiError struct {
	Context       string
	Message       string
	ExceptionName string
}

type requestError struct {
	Code    int
	Message string
}

func (r requestError) Error() string {
	return r.Message
}

type getRepositoriesResponse struct {
	pagedResponse
	Values []repository
}

type bitClient struct {
	sling *sling.Sling
	Path  string
}

// List is an interface for adding items to a list
type List interface {
	Push(string)
}

// DiffItem is a diff struct for an inidividual file
type DiffItem struct {
	raw    string
	fPath  string
	commit string
}

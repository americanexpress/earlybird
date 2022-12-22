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
	"strings"
)

const (
	pathSep = " b/"
)

// GetHashKey returns the hash key identifier for the diff
func (d *DiffItem) GetHashKey() string {
	if d.commit != "" {
		return fmt.Sprintf("%s:%s", d.commit, d.fPath)
	}
	return d.fPath
}

// Diff is a list of split diffs
type Diff struct {
	Items     []DiffItem
	Error     error
	commitTmp string
}

// Push a diff on to the list
func (d *Diff) Push(s string) {

	var commitHeader, commit string

	if beginsWithHash(s) {
		commitHeader, s = split(s, "\n")
		commit = extractHash(commitHeader)
		d.commitTmp = commit
	}
	// add commit to diffs within each diff which do not
	// have the commit on the line directly above them
	if d.commitTmp != "" && commit == "" {
		commit = d.commitTmp
	}

	fPath, err := extractFilePath(s)
	if err != nil {
		d.Error = err
		return
	}

	d.Items = append(d.Items, DiffItem{
		raw:    s,
		fPath:  fPath,
		commit: commit,
	})
}

// split out logic from scan.go
func beginsWithHash(s string) bool {
	return strings.Contains(s, "commit")
}

func split(s, sep string) (string, string) {
	// Empty string should just return empty
	if len(s) == 0 {
		return s, s
	}
	slice := strings.SplitN(s, sep, 2)
	// Incase no separator was present
	if len(slice) == 1 {
		return slice[0], ""
	}
	return slice[0], slice[1]
}

func extractFilePath(in string) (string, error) {
	pathBIndex := strings.Index(in, pathSep)
	newLineIndex := strings.Index(in, "\n")
	if pathBIndex >= 0 && newLineIndex > pathBIndex {
		return in[pathBIndex+len(pathSep) : newLineIndex], nil
	}
	return "", fmt.Errorf("Not valid diff content:\n%s", in)
}

func extractHash(in string) string {
	if len(strings.Fields(in)) >= 2 {
		return strings.Fields(in)[1]
	}
	return ""
}

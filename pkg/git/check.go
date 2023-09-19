/*
 * Copyright 2021 American Express
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
	"io"
	"strings"

	"github.com/americanexpress/earlybird/v4/pkg/scan"
)

//ParseGitLog parses the git log into the earlybird file format
func ParseGitLog(r io.Reader) (fileList []scan.File, err error) {
	diff := Diff{}
	err = splitDiffs(r, &diff)
	if err != nil {
		return nil, err
	}

	for _, d := range diff.Items {
		//Build file here
		curFile := scan.File{
			Name: "buffer",
			Path: fmt.Sprintf("%s:%s", d.commit, d.fPath),
		}
		//Append lines
		var line scan.Line
		lineNum := 1
		for _, lineText := range strings.Split(strings.TrimSuffix(d.raw, "\n"), "\n") {
			line.LineNum = lineNum
			line.LineValue = lineText
			line.FilePath = curFile.Path
			lineNum++
			curFile.Lines = append(curFile.Lines, line)
		}
		//Append commit file
		fileList = append(fileList, curFile)
	}

	return fileList, nil
}

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

import "bytes"

const (
	diffSep = "diff --git a"
)

func isDiffHeader(data *[]byte, index int) bool {
	diffSepLen := len(diffSep)
	diffSepEndIndex := index + diffSepLen
	dataLen := len(*data) - 1
	return diffSepEndIndex < dataLen && string((*data)[index:diffSepEndIndex]) == diffSep
}
func getPreviousLineIndex(data []byte, index int) int {
	i := bytes.LastIndex((data)[0:index], []byte("\n"))
	// if no previous line then must be on the very first line
	if i < 0 {
		return 0
	}
	return i
}

// ScanDiffs splits on the diff of an inidividual file
func ScanDiffs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	k, newLineIndex := 0, 0
	dataLen := len(data) - 1

	// loop until no more bytes left to read in this chunk of data
	for k < dataLen {
		i := bytes.IndexByte(data[k:], '\n')
		if i < 0 {
			k = dataLen
			continue
		}

		// k = index of scanned through data so far
		// i = index after last \n char
		// 1 = start at next byte
		newLineIndex = k + i + 1

		if beginsWithHash(string(data[newLineIndex:])) {
			return newLineIndex, dropCR(data[0 : newLineIndex-1]), nil
		}

		if isDiffHeader(&data, newLineIndex) {
			// if previous line does not begin with a hash then separate on diff headers
			prevLineIndex := getPreviousLineIndex(data, newLineIndex-1)
			if !beginsWithHash(string(data[prevLineIndex : k+i])) {
				return newLineIndex, dropCR(data[0 : newLineIndex-1]), nil
			}
		}

		// keep advancing through data
		// k = index of scanned so far
		// i = index of new line
		// 1 = start at next byte after previous new line
		k += i + 1
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

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

package configupdate

import (
	"os"
	"testing"
)

func Test_downloadFile(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}
	path, _ := os.Getwd()
	path += "test.txt"
	err := downloadFile(path, "http://httpstat.us/200") //Check if page loads
	if err != nil {
		t.Errorf("Failed to download file: %v", err)
	} else {
		os.Remove(path) //Delete downloaded file
	}
}

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

package file

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_ExtendedGitIgnoreSamples(t *testing.T) {
	// Define test cases for each language/sample
	type testCase struct {
		filePath      string // Path relative to searchDir
		shouldIgnore  bool
	}

	samples := []struct {
		name       string
		ignoreFile string
		cases      []testCase
	}{
		{
			name:       "Go",
			ignoreFile: "test_data/gitignore_samples/Go.gitignore",
			cases: []testCase{
				// Ignored patterns
				{"myprogram.exe", true},
				{"test.exe~", true},
				{"pkg.dll", true},
				{"build/main.test", true},
				{"coverage.out", true},
				{"profile.cov", true},
				{"go.work", true},
				{".env", true},
				{"config/.env", true}, // .env can be anywhere (depends on matcher root, assuming root for now)
				
				// Not ignored
				{"main.go", false},
				{"go.mod", false},
				{"readme.md", false},
				{"pkg/main.go", false},
			},
		},
		{
			name:       "Python",
			ignoreFile: "test_data/gitignore_samples/Python.gitignore",
			cases: []testCase{
				// Ignored patterns
				{"__pycache__/cache.pyc", true},
				{"src/__pycache__/cache.pyc", true},
				{"module.so", true},
				{"build/lib/pkg", true},
				{"dist/package-1.0.tar.gz", true},
				{".env", true},
				{".venv/bin/activate", true},
				{".idea/workspace.xml", false}, // Python.gitignore doesn't ignore .idea by default (usually global)
				{"htmlcov/index.html", true},
				
				// Not ignored
				{"main.py", false},
				{"setup.py", false},
				{"requirements.txt", false},
				{"src/module.py", false},
			},
		},
		{
			name:       "Node",
			ignoreFile: "test_data/gitignore_samples/Node.gitignore",
			cases: []testCase{
				// Ignored patterns
				{"node_modules/package.json", true},
				{"logs/debug.log", true},
				{"npm-debug.log", true},
				{"coverage/lcov.info", true},
				{".env", true},
				{".env.local", true},
				{"dist/app.js", true},
				{".DS_Store", false}, // standard macOS ignore not in Node.gitignore usually (but good to check it's not falsely positive)
				
				// Negation checks (previously unsupported)
				{".env.example", false}, // Explicitly un-ignored: !.env.example
				
				// Not ignored
				{"package.json", false},
				{"src/index.js", false},
				{"public/index.html", false},
			},
		},
		{
			name:       "Java",
			ignoreFile: "test_data/gitignore_samples/Java.gitignore",
			cases: []testCase{
				// Ignored patterns
				{"Main.class", true},
				{"build/classes/Shape.class", true},
				{"app.jar", true},
				{"lib/dependency.war", true},
				{"server.log", true},
				{"hs_err_pid1234.log", true},
				
				// Not ignored
				{"Main.java", false},
				{"gradlew", false},
				{"pom.xml", false},
			},
		},
	}

	for _, sample := range samples {
		t.Run(sample.name, func(t *testing.T) {
			// Mock the global ignorePatterns with the sample file
			// We use empty string for the first arg (searchDir) because getIgnorePatterns joins it with .ge_ignore
			// Here we are testing the "ignoreFile" argument explicitly or we can just pass the path as the "ignoreFile" arg.
			// The getIgnorePatterns signature is: (filePath, ignoreFile string, verbose bool)
			// It loads .ge_ignore from filePath AND the specific ignoreFile.
			
			// To test JUST the sample file, we can pass a dummy path for filePath
			absPath, _ := filepath.Abs(sample.ignoreFile)
			_, matcher := getIgnorePatterns("/tmp/nonexistent", absPath, false)
			
			for _, tc := range sample.cases {
				// go-git matcher expects path components
				// We also need to decide isDir. For simple tests, we assume path ending in / isDir ??
				// But the inputs in test cases don't always end in /. 
				// However, if the pattern matched a directory, usually convention is it matches contents.
				
				// Quick path splitter
				// Clean path first
				cleanPath := tc.filePath
				parts := []string{}
				for _, p := range strings.Split(cleanPath, "/") {
					if p != "" {
						parts = append(parts, p)
					}
				}
				
				// Guess isDir?
				isDir := strings.HasSuffix(tc.filePath, "/")
				// If expected to be ignored and it looks like a directory pattern in gitignore...
				// But we are testing if a FILE at that path is ignored.
				// In gitignore validation, usually we verify files.
				
				got := matcher.Match(parts, isDir)
				if got != tc.shouldIgnore {
					t.Errorf("[%s] File '%s': expected ignore=%v, got=%v", sample.name, tc.filePath, tc.shouldIgnore, got)
				}
			}
		})
	}
}

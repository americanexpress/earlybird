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

package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/americanexpress/earlybird/v4/pkg/scan"
)

func sampleHit(filename string, line int, lineVal, matchVal, caption string, code int) scan.Hit {
	return scan.Hit{
		Filename:   filename,
		Line:       line,
		LineValue:  lineVal,
		MatchValue: matchVal,
		Caption:    caption,
		Category:   "password-secret",
		Code:       code,
		Severity:   "high",
		Confidence: "high",
	}
}

// TestStableFindingFingerprints verifies that fingerprints remain constant even if
// line numbers shift due to unrelated code changes, or whitespace is padded.
func TestStableFindingFingerprints(t *testing.T) {
	hitOriginal := sampleHit("pkg/auth/login.go", 42, "password = \"secret123\"", "secret123", "Hardcoded Password", 101)
	hitShiftedLine := sampleHit("pkg/auth/login.go", 188, "password = \"secret123\"", "secret123", "Hardcoded Password", 101)
	hitIndented := sampleHit("pkg/auth/login.go", 42, "    password = \"secret123\"    ", "secret123", "Hardcoded Password", 101)

	fp1 := ComputeFingerprint(hitOriginal)
	fp2 := ComputeFingerprint(hitShiftedLine)
	fp3 := ComputeFingerprint(hitIndented)

	if fp1 != fp2 {
		t.Errorf("Fingerprints differ after line shift: %s vs %s", fp1, fp2)
	}
	if fp1 != fp3 {
		t.Errorf("Fingerprints differ after whitespace padding: %s vs %s", fp1, fp3)
	}

	// Different secrets must produce different fingerprints
	hitDifferentSecret := sampleHit("pkg/auth/login.go", 42, "password = \"another_secret\"", "another_secret", "Hardcoded Password", 101)
	fpDiff := ComputeFingerprint(hitDifferentSecret)
	if fp1 == fpDiff {
		t.Errorf("Different secrets should not have matching fingerprints")
	}
}

// TestExistingBaselineFindings verifies that findings in the baseline are recognized and ignored.
func TestExistingBaselineFindings(t *testing.T) {
	b := New()
	hit1 := sampleHit("src/main.py", 10, "api_key = 'xyz123'", "xyz123", "API Key", 201)
	hit2 := sampleHit("src/db.go", 25, "db_pass = 'pass123'", "pass123", "DB Password", 202)

	b.Add(hit1)
	b.Add(hit2)

	if !b.IsIgnored(hit1) {
		t.Errorf("Expected hit1 to be ignored as it is in the baseline")
	}
	if !b.IsIgnored(hit2) {
		t.Errorf("Expected hit2 to be ignored as it is in the baseline")
	}
}

// TestNewFindings verifies that new findings not present in the baseline are NOT ignored.
func TestNewFindings(t *testing.T) {
	b := New()
	baselineHit := sampleHit("src/main.py", 10, "api_key = 'xyz123'", "xyz123", "API Key", 201)
	b.Add(baselineHit)

	newHit := sampleHit("src/new_feature.py", 5, "token = 'new_token_456'", "new_token_456", "Auth Token", 203)

	if b.IsIgnored(newHit) {
		t.Errorf("New finding should NOT be ignored by baseline")
	}

	hits := []scan.Hit{baselineHit, newHit}
	filtered := FilterHits(hits, b)
	if len(filtered) != 1 {
		t.Fatalf("Expected 1 new hit after filtering, got %d", len(filtered))
	}
	if filtered[0].MatchValue != "new_token_456" {
		t.Errorf("Expected new_token_456, got %s", filtered[0].MatchValue)
	}
}

// TestRemovedFindings verifies that if code fixes/removes a finding from baseline,
// the scanner continues normally without errors.
func TestRemovedFindings(t *testing.T) {
	b := New()
	fixedHit := sampleHit("src/legacy.go", 10, "old_key = 'fixed'", "fixed", "API Key", 201)
	remainingHit := sampleHit("src/db.go", 20, "db_key = 'still_here'", "still_here", "DB Key", 202)

	b.Add(fixedHit)
	b.Add(remainingHit)

	// Scan only produces remainingHit now (fixedHit was removed from code)
	scannedHits := []scan.Hit{remainingHit}
	filtered := FilterHits(scannedHits, b)

	if len(filtered) != 0 {
		t.Errorf("Expected 0 new hits after filtering fixed baseline, got %d", len(filtered))
	}
}

// TestDuplicateFindings verifies that duplicates within a baseline or added multiple times
// are safely deduplicated.
func TestDuplicateFindings(t *testing.T) {
	b := New()
	hit := sampleHit("src/app.go", 15, "secret = 'abc'", "abc", "Secret", 301)

	added1 := b.Add(hit)
	added2 := b.Add(hit)

	if !added1 {
		t.Errorf("First add should return true")
	}
	if added2 {
		t.Errorf("Second add should return false (duplicate)")
	}
	if b.Count() != 1 {
		t.Errorf("Expected baseline count 1, got %d", b.Count())
	}

	// Slice deduplication
	dupes := []scan.Hit{hit, hit, hit}
	unique := DeduplicateHits(dupes)
	if len(unique) != 1 {
		t.Errorf("Expected 1 unique hit, got %d", len(unique))
	}
}

// TestInvalidBaselineFiles verifies error handling for missing, corrupted, and empty files.
func TestInvalidBaselineFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Missing file
	_, err := Load(filepath.Join(tempDir, "non_existent.json"))
	if err == nil {
		t.Errorf("Expected error loading non-existent baseline file")
	}

	// 2. Empty path
	_, err = Load("")
	if err == nil {
		t.Errorf("Expected error loading empty path")
	}

	// 3. Empty file
	emptyFile := filepath.Join(tempDir, "empty.json")
	if err := os.WriteFile(emptyFile, []byte("   \n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = Load(emptyFile)
	if err != ErrEmptyBaselineFile {
		t.Errorf("Expected ErrEmptyBaselineFile, got %v", err)
	}

	// 4. Corrupt JSON
	corruptFile := filepath.Join(tempDir, "corrupt.json")
	if err := os.WriteFile(corruptFile, []byte("{ not valid json !!!"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = Load(corruptFile)
	if err == nil {
		t.Errorf("Expected error loading corrupt JSON file")
	}
}

// TestChangedFilePaths verifies that when a file is moved to another directory or
// Windows/Unix path separators differ, the baseline still matches.
func TestChangedFilePaths(t *testing.T) {
	b := New()
	hitOriginal := sampleHit("src/services/auth.go", 10, "token = 'token_abc'", "token_abc", "Auth Token", 401)
	b.Add(hitOriginal)

	// Same file with Windows separator
	hitWindowsPath := sampleHit("src\\services\\auth.go", 10, "token = 'token_abc'", "token_abc", "Auth Token", 401)
	if !b.IsIgnored(hitWindowsPath) {
		t.Errorf("Expected Windows-path finding to be recognized by baseline")
	}

	// File moved to new location (fallback content fingerprint)
	hitMovedFile := sampleHit("pkg/internal/auth.go", 10, "token = 'token_abc'", "token_abc", "Auth Token", 401)
	if !b.IsIgnored(hitMovedFile) {
		t.Errorf("Expected relocated finding with same base filename to be recognized by content fingerprint")
	}
}

// TestSaveAndLoadRoundtrip verifies saving findings to disk and loading them back.
func TestSaveAndLoadRoundtrip(t *testing.T) {
	tempDir := t.TempDir()
	baselinePath := filepath.Join(tempDir, "sub", "baseline.json")

	hits := []scan.Hit{
		sampleHit("pkg/a.go", 1, "sec = '1'", "1", "Sec1", 1),
		sampleHit("pkg/b.go", 2, "sec = '2'", "2", "Sec2", 2),
		sampleHit("pkg/a.go", 1, "sec = '1'", "1", "Sec1", 1), // dupe
	}

	err := Save(baselinePath, hits, "v4.0.0")
	if err != nil {
		t.Fatalf("Failed to save baseline: %v", err)
	}

	loaded, err := Load(baselinePath)
	if err != nil {
		t.Fatalf("Failed to load baseline: %v", err)
	}

	if loaded.Count() != 2 {
		t.Errorf("Expected 2 unique hits loaded, got %d", loaded.Count())
	}
	if !loaded.IsIgnored(hits[0]) {
		t.Errorf("Expected hits[0] to be ignored")
	}
	if !loaded.IsIgnored(hits[1]) {
		t.Errorf("Expected hits[1] to be ignored")
	}
}

// TestLoadFromRawArrayFormat verifies that a baseline file formatted as a raw JSON array
// of hits is parsed correctly.
func TestLoadFromRawArrayFormat(t *testing.T) {
	tempDir := t.TempDir()
	rawPath := filepath.Join(tempDir, "raw_baseline.json")

	hits := []scan.Hit{
		sampleHit("file1.go", 5, "pass = 'raw1'", "raw1", "Pass", 1),
	}

	bytes, _ := json.Marshal(hits)
	if err := os.WriteFile(rawPath, bytes, 0644); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(rawPath)
	if err != nil {
		t.Fatalf("Failed to load raw array baseline: %v", err)
	}
	if loaded.Count() != 1 {
		t.Errorf("Expected 1 hit loaded from raw array, got %d", loaded.Count())
	}
	if !loaded.IsIgnored(hits[0]) {
		t.Errorf("Expected hits[0] to be ignored")
	}
}

// TestBackwardCompatibility verifies that nil baseline gracefully passes through all hits.
func TestBackwardCompatibility(t *testing.T) {
	var nilBaseline *Baseline
	hit := sampleHit("any.go", 1, "key = 'val'", "val", "Key", 1)

	if nilBaseline.IsIgnored(hit) {
		t.Errorf("Nil baseline should never ignore findings")
	}

	hits := []scan.Hit{hit}
	filtered := FilterHits(hits, nilBaseline)
	if len(filtered) != 1 {
		t.Errorf("Nil baseline should return original hits unfiltered")
	}
}

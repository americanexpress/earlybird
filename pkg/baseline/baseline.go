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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/americanexpress/earlybird/v4/pkg/scan"
)

var (
	// ErrEmptyBaselineFile is returned when the baseline file is empty
	ErrEmptyBaselineFile = errors.New("baseline file is empty")
	// ErrInvalidBaselineFormat is returned when the baseline file cannot be parsed
	ErrInvalidBaselineFormat = errors.New("invalid baseline file format")
)

// BaselineFile represents the serialized baseline format
type BaselineFile struct {
	Version     string     `json:"version,omitempty"`
	GeneratedAt string     `json:"generated_at,omitempty"`
	HitCount    int        `json:"hit_count"`
	Hits        []scan.Hit `json:"hits"`
}

// Baseline manages a set of known findings to exclude during security scans
type Baseline struct {
	fingerprints        map[string]bool
	contentFingerprints map[string]bool
	hits                []scan.Hit
}

// New creates an empty Baseline instance
func New() *Baseline {
	return &Baseline{
		fingerprints:        make(map[string]bool),
		contentFingerprints: make(map[string]bool),
		hits:                make([]scan.Hit, 0),
	}
}

// NormalizePath standardizes file paths across operating systems
func NormalizePath(path string) string {
	cleaned := filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))
	cleaned = strings.TrimPrefix(cleaned, "./")
	return cleaned
}

// ComputeFingerprint produces a stable cryptographic digest of a finding.
// Line numbers are intentionally excluded so findings remain stable when
// lines are added or removed elsewhere in the file.
func ComputeFingerprint(hit scan.Hit) string {
	normalizedPath := NormalizePath(hit.Filename)
	trimmedLine := strings.TrimSpace(hit.LineValue)
	trimmedMatch := strings.TrimSpace(hit.MatchValue)

	raw := fmt.Sprintf("%s|%d|%s|%s|%s|%s",
		normalizedPath,
		hit.Code,
		hit.Caption,
		hit.Category,
		trimmedLine,
		trimmedMatch,
	)

	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// ComputeContentFingerprint produces a path-agnostic content fingerprint based on the
// base filename and match details. This enables matching findings even if a file
// was relocated within the project tree.
func ComputeContentFingerprint(hit scan.Hit) string {
	baseName := filepath.Base(hit.Filename)
	trimmedLine := strings.TrimSpace(hit.LineValue)
	trimmedMatch := strings.TrimSpace(hit.MatchValue)

	raw := fmt.Sprintf("%s|%d|%s|%s|%s|%s",
		baseName,
		hit.Code,
		hit.Caption,
		hit.Category,
		trimmedLine,
		trimmedMatch,
	)

	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// Add adds a finding to the baseline, ignoring duplicates
func (b *Baseline) Add(hit scan.Hit) bool {
	fp := ComputeFingerprint(hit)
	if b.fingerprints[fp] {
		return false
	}

	b.fingerprints[fp] = true
	b.contentFingerprints[ComputeContentFingerprint(hit)] = true
	b.hits = append(b.hits, hit)
	return true
}

// IsIgnored reports whether the given finding is present in the baseline
func (b *Baseline) IsIgnored(hit scan.Hit) bool {
	if b == nil {
		return false
	}

	// 1. Check primary exact fingerprint
	fp := ComputeFingerprint(hit)
	if b.fingerprints[fp] {
		return true
	}

	// 2. Check fallback content fingerprint (for changed relative paths / moved files)
	cfp := ComputeContentFingerprint(hit)
	if b.contentFingerprints[cfp] {
		return true
	}

	return false
}

// Hits returns all unique baseline findings
func (b *Baseline) Hits() []scan.Hit {
	return b.hits
}

// Count returns the number of findings in the baseline
func (b *Baseline) Count() int {
	return len(b.hits)
}

// Load loads a baseline file from disk. It accepts both standard BaselineFile / Report
// JSON formats (with a "hits" field) and raw JSON arrays of scan.Hit.
func Load(filePath string) (*Baseline, error) {
	if filePath == "" {
		return nil, errors.New("baseline file path cannot be empty")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read baseline file: %w", err)
	}

	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return nil, ErrEmptyBaselineFile
	}

	b := New()

	// Try unmarshaling as BaselineFile / Report format first
	var fileFormat BaselineFile
	if err := json.Unmarshal(data, &fileFormat); err == nil && fileFormat.Hits != nil {
		for _, hit := range fileFormat.Hits {
			b.Add(hit)
		}
		return b, nil
	}

	// Fallback to raw array of hits
	var rawHits []scan.Hit
	if err := json.Unmarshal(data, &rawHits); err == nil {
		for _, hit := range rawHits {
			b.Add(hit)
		}
		return b, nil
	}

	return nil, fmt.Errorf("%w: failed to parse JSON in %s", ErrInvalidBaselineFormat, filePath)
}

// Save writes a slice of findings to a baseline file on disk.
func Save(filePath string, hits []scan.Hit, version string) error {
	if filePath == "" {
		return errors.New("baseline output path cannot be empty")
	}

	b := New()
	for _, hit := range hits {
		b.Add(hit)
	}

	baselineFile := BaselineFile{
		Version:     version,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		HitCount:    b.Count(),
		Hits:        b.Hits(),
	}

	bytes, err := json.MarshalIndent(baselineFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal baseline data: %w", err)
	}

	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for baseline file: %w", err)
		}
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write baseline file: %w", err)
	}

	return nil
}

// FilterHits filters out any hits present in the baseline, returning only new findings.
func FilterHits(hits []scan.Hit, b *Baseline) []scan.Hit {
	if b == nil {
		return hits
	}

	var newHits []scan.Hit
	for _, hit := range hits {
		if !b.IsIgnored(hit) {
			newHits = append(newHits, hit)
		}
	}
	return newHits
}

// DeduplicateHits removes duplicate hits using stable fingerprinting.
func DeduplicateHits(hits []scan.Hit) []scan.Hit {
	seen := make(map[string]bool)
	var unique []scan.Hit
	for _, hit := range hits {
		fp := ComputeFingerprint(hit)
		if !seen[fp] {
			seen[fp] = true
			unique = append(unique, hit)
		}
	}
	return unique
}

// Unused import guard
var _ = strconv.Itoa

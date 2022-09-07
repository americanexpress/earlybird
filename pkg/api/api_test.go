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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/scan"
	"github.com/americanexpress/earlybird/pkg/utils"
)

var cfg = cfgreader.EarlybirdConfig{
	ConfigDir:               path.Join(utils.MustGetWD(), utils.GetConfigDir()),
	LabelsConfigDir:         path.Join(utils.MustGetWD(), utils.GetConfigDir(), "labels"),
	FalsePositivesConfigDir: path.Join(utils.MustGetWD(), utils.GetConfigDir(), "falsepositives"),
	RulesConfigDir:          path.Join(utils.MustGetWD(), utils.GetConfigDir(), "rules"),
	IgnoreFile:              path.Join(utils.MustGetWD(), "../../.ge_ignore"),
	SeverityDisplayLevel:    4,
	SeverityFailLevel:       4,
	ConfidenceDisplayLevel:  4,
	ConfidenceFailLevel:     4,
	MaxFileSize:             10240000,
	WorkLength:              2500,
	EnabledModulesMap:       map[string]string{"content": "content.yaml", "password-secret": "password-secret.yaml"},
	WorkerCount:             100,
}

func init() {
	scan.Init(cfg)
}

func TestScan(t *testing.T) {
	fmt.Println("Testing Scan")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("scan", "sample.py")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write([]byte(`password="SampleFinding678#"`))
	if err != nil {
		t.Fatal(err)
	}

	//Submit example file containing a secret
	req, err := http.NewRequest("POST", "/scan", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	writer.Close()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Scan(cfg))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Scan returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse body to verify secrets hits on URL
	var results scan.Report
	err = json.NewDecoder(rr.Body).Decode(&results)
	if err != nil {
		t.Errorf("Failed to parse result from Scan: %v", err)
	}
	if len(results.Hits) == 0 {
		t.Errorf("Failed to locate secrets in example file: %v", results)
	}
}

func TestGITScan(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}
	fmt.Println("Testing GITScan")
	//Submit example github repository containing secrets
	req, err := http.NewRequest("GET", "/scan/git?url=https://github.com/carnal0wnage/fake_commited_secrets", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GITScan(cfg))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GITScan returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse body to verify secrets hits on URL
	var results scan.Report
	err = json.NewDecoder(rr.Body).Decode(&results)
	if err != nil {
		t.Errorf("Failed to parse result from GITScan: %v", err)
	}
	if len(results.Hits) == 0 {
		t.Errorf("Failed to locate secrets in example repository: %v", results)
	}
}

func TestLabels(t *testing.T) {
	fmt.Println("Testing Labels")
	//fetch all labels per category
	req, err := http.NewRequest("GET", "/labels", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Labels(cfg.Version, scan.Labels))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Labels returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse body to verify secrets hits on URL
	var results LabelsResponse
	err = json.NewDecoder(rr.Body).Decode(&results)
	if err != nil {
		t.Errorf("Failed to parse result from Labels: %v", err)
	}
	if len(results.Labels) == 0 {
		t.Errorf("Failed to load Labels: %v", results)
	}
}

func TestCategories(t *testing.T) {
	fmt.Println("Testing Categorizes")
	//fetch all labels per category
	req, err := http.NewRequest("GET", "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Categories(cfg.Version, scan.CombinedRules))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Categories returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse body to verify secrets hits on URL
	var results CategoriesResponse
	err = json.NewDecoder(rr.Body).Decode(&results)
	if err != nil {
		t.Errorf("Failed to parse result from Categories: %v", err)
	}
	if len(results.Categories) == 0 {
		t.Errorf("Failed to load categories: %v", results)
	}
}

func TestLabelsPerCategory(t *testing.T) {
	fmt.Println("Testing LabelsPerCategory")
	//fetch all labels per category
	req, err := http.NewRequest("GET", "/categorylabels", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LabelsPerCategory(cfg.Version, scan.Labels))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("LabelsPerCategory returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse body to verify secrets hits on URL
	var results CategoryLabelsResponse
	err = json.NewDecoder(rr.Body).Decode(&results)
	if err != nil {
		t.Errorf("Failed to parse result from LabelsPerCategory: %v", err)
	}
	if len(results.CategoryLabels) == 0 {
		t.Errorf("Failed to load Labels per Category: %v", results)
	}
}

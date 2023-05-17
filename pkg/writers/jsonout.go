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

package writers

import (
	"encoding/json"
	"fmt"
	cfgReader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/file"
	"github.com/americanexpress/earlybird/pkg/scan"
	"io/ioutil"
	"os"
	"time"
)

// WriteJSON takes the hits, converts them into JSON report and passing report to reportToJSONWriter().
func WriteJSON(hits <-chan scan.Hit, config cfgReader.EarlybirdConfig, fileContext file.Context, fileName string) (err error) {
	start := time.Now()
	var Hits []scan.Hit
	for hit := range hits {
		Hits = append(Hits, hit)
	}

	report := scan.Report{
		Hits:          Hits,
		HitCount:      len(Hits),
		Skipped:       fileContext.SkippedFiles,
		Ignore:        fileContext.IgnorePatterns,
		Version:       config.Version,
		Modules:       config.EnabledModules,
		Threshold:     config.SeverityDisplayLevel,
		FilesScanned:  len(fileContext.Files),
		RulesObserved: len(scan.CombinedRules),
		StartTime:     start.UTC().Format(time.RFC3339),
		EndTime:       time.Now().UTC().Format(time.RFC3339),
		Duration:      fmt.Sprintf("%d ms", time.Since(start)/time.Millisecond),
	}
	_, err = reportToJSONWriter(report, fileName)
	return err
}

//reportToJSONWriter Outputs an object as a JSON blob to an output file or console
func reportToJSONWriter(v interface{}, fileName string) (s string, err error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return "", err
	}
	if fileName == "" {
		_, err = os.Stdout.Write(b)
	} else {
		err = ioutil.WriteFile(fileName, b, 0666)
	}
	return string(b), err

}

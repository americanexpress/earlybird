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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	cfgreader "github.com/americanexpress/earlybird/v4/pkg/config"
)

//UpdateConfigFiles updates all of the modules via the defined module URL in earlybird.json
func UpdateConfigFiles(configDir, rulesConfigDir, appConfigPath, appConfigURL string, ruleModulesFilenameMap map[string]string) error {
	for _, fileName := range ruleModulesFilenameMap {
		moduleFilePath := path.Join(rulesConfigDir, fileName)
		log.Println("Updating ", moduleFilePath)

		parsedUrl, err := url.Parse(cfgreader.Settings.ConfigBaseUrl)

		if err != nil {
			return err
		}

		parsedUrl.Path = path.Join(parsedUrl.Path, fileName)

		err = downloadFile(moduleFilePath, parsedUrl.String())

		if err != nil {
			return err
		}
	}

	log.Println("Updating ", appConfigPath)
	return downloadFile(path.Join(configDir, "earlybird.json"), appConfigURL)
}

func downloadFile(path string, url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("downloading file from %s: %v", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non 200 status code: status=%d, response=%v", resp.StatusCode, string(b))
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	err = ioutil.WriteFile(path, b, 0666)
	if err != nil {
		return fmt.Errorf("writing file at %s: %v", path, err)
	}
	return nil
}

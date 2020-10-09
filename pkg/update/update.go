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

package configupdate

import (
	"fmt"
	"io/ioutil"
	"net/http"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
)

//UpdateConfigFiles updates all of the modules via the defined module URL in earlybird.json
func UpdateConfigFiles(configDir string, appConfigPath string, appConfigURL string) error {
	var moduleFilePath string
	for _, module := range cfgreader.Settings.ModuleConfigs {
		moduleFilePath = configDir + module.Name + ".json"
		fmt.Println("Updating ", moduleFilePath)
		if module.ConfigURL != "" {
			err := downloadFile(moduleFilePath, module.ConfigURL)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Updating ", appConfigPath)
	return downloadFile(moduleFilePath, appConfigURL)
}

func downloadFile(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading file from %s: %v", url, err)
	}
	defer resp.Body.Close()
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

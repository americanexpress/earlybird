/*
 * Copyright 2023 American Express
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

package cfgreader

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"unicode"

	"github.com/ghodss/yaml"
)

// Settings contains imported earlybird.json configuration
var Settings Configs

// LoadConfig parses json configuration file into structure
func LoadConfig(cfg interface{}, path string) (err error) {
	file, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer file.Close()

	// read our opened json file as a byte array.
	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return readErr
	}

	byteValue, err := toJson(data)
	if err != nil {
		return err
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into arrays which we defined above
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err
	}

	return err
}

// toJson returns json from data. If data is already json this fn is a noop, if data is yaml
// the yaml is converted to json
func toJson(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}

	return yaml.YAMLToJSON(data)
}

// hasJSONPrefix determines if data has an opening curly brace, signifying the file is json
func hasJSONPrefix(data []byte) bool {
	jsonPrefix := []byte("{")
	return hasPrefix(data, jsonPrefix)
}

// hasPrefix determines if the bytes has a particular prefix
func hasPrefix(data, prefix []byte) bool {
	trimmed := bytes.TrimLeftFunc(data, unicode.IsSpace)

	return bytes.HasPrefix(trimmed, prefix)
}

// TranslateLevelID returns the text value of a level integer (e.g., 2 --> high)
func (cfg *Configs) TranslateLevelID(level int) string {
	for _, levelConfig := range cfg.LevelConfigs {
		if levelConfig.ID == level {
			return levelConfig.Name
		}
	}
	return "low"
}

// TranslateLevelName returns the int level of a string value (e.g., high --> 2)
func (cfg *Configs) TranslateLevelName(level string) int {
	for _, levelConfig := range cfg.LevelConfigs {
		if levelConfig.Name == level {
			return levelConfig.ID
		}
	}
	return 4
}

// GetLevelNames returns a slice of strings for the level names
func (cfg *Configs) GetLevelNames() (levels []string) {
	for _, levelConfig := range cfg.LevelConfigs {
		levels = append(levels, levelConfig.Name)
	}
	return levels
}

// GetLevelMap returns a map of severity levels
func (cfg *Configs) GetLevelMap() (levelMap map[string]int) {
	levelMap = make(map[string]int)
	for _, levelConfig := range cfg.LevelConfigs {
		levelMap[levelConfig.Name] = levelConfig.ID
	}
	return levelMap
}

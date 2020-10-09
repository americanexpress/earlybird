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

package cfgreader

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"os"

	"github.com/americanexpress/earlybird/pkg/utils"
)

//Settings contains imported earlybird.json configuration
var Settings Configs

func init() {
	configPath := utils.GetConfigDir() + "earlybird.json"
	err := LoadConfig(&Settings, configPath)
	if err != nil {
		log.Fatal("Failed to load Earlybird config", err)
	}
}

//LoadConfig parses json configuration file into structure
func LoadConfig(cfg interface{}, path string) (err error) {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened json file as a byte array.
	byteValue, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
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

//TranslateLevelID returns the text value of a level integer (e.g., 2 --> high)
func (cfg *Configs) TranslateLevelID(level int) string {
	for _, levelConfig := range cfg.LevelConfigs {
		if levelConfig.ID == level {
			return levelConfig.Name
		}
	}
	return "low"
}

//TranslateLevelName returns the int level of a string value (e.g., high --> 2)
func (cfg *Configs) TranslateLevelName(level string) int {
	for _, levelConfig := range cfg.LevelConfigs {
		if levelConfig.Name == level {
			return levelConfig.ID
		}
	}
	return 4
}

//GetLevelNames returns a slice of strings for the level names
func (cfg *Configs) GetLevelNames() (levels []string) {
	for _, levelConfig := range cfg.LevelConfigs {
		levels = append(levels, levelConfig.Name)
	}
	return levels
}

//GetLevelMap returns a map of severity levels
func (cfg *Configs) GetLevelMap() (levelMap map[string]int) {
	levelMap = make(map[string]int)
	for _, levelConfig := range cfg.LevelConfigs {
		levelMap[levelConfig.Name] = levelConfig.ID
	}
	return levelMap
}

//GetAvailableModules returns all Earlybird available modules
func (cfg *Configs) GetAvailableModules() (modules []string) {
	for _, module := range cfg.ModuleConfigs {
		modules = append(modules, module.Name)
	}
	return modules
}

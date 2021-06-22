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

package cfgreader

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/americanexpress/earlybird/pkg/utils"
)

var (
	config     Configs
	configPath = utils.GetConfigDir() + "earlybird.json"
)

func init() {
	if err := LoadConfig(&config, configPath); err != nil {
		fmt.Println("LoadConfig() = error:", err, " want error nil")
	}
}

func TestLoadConfig(t *testing.T) {
	var settings Configs
	if err := LoadConfig(&settings, configPath); err != nil && len(settings.Version) == 0 {
		t.Errorf("LoadConfigs() = %v, want non nil value, loaded empty config", settings)
	}
}

func TestTranslateLevelID(t *testing.T) {
	result := "critical"
	if got := config.TranslateLevelID(1); got != result {
		t.Errorf("TranslateLevelID() = %v, want %v", got, result)
	}
}

func TestTranslateLevelName(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Convert name to int",
			args: args{
				level: "high",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.TranslateLevelName(tt.args.level); got != tt.want {
				t.Errorf("TranslateLevelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLevelNames(t *testing.T) {
	if gotLevels := config.GetLevelNames(); len(gotLevels) == 0 {
		t.Errorf("GetLevelNames() = %v, wanted non empty array", gotLevels)
	}
}

func TestGetLevelMap(t *testing.T) {
	gotLevelMap := config.GetLevelMap()
	if _, ok := gotLevelMap["low"]; !ok {
		t.Errorf("GetLevelMap() = %v, want non empty map", gotLevelMap)
	}
}

func TestHasJSONPrefix(t *testing.T) {
	cases := map[string]struct {
		expected bool
		input    []byte
	}{
		"input is json": {
			expected: true,
			input:    []byte("{\"this\": \"is JSON\"}"),
		},
		"input is not json": {
			expected: false,
			input:    []byte("this is not json"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if got := hasJSONPrefix(tc.input); tc.expected != got {
				t.Fatalf("Did not get expected result, got: %t, want: %t", got, tc.expected)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	cases := map[string]struct {
		expected []byte
		input    []byte
	}{
		"input is json": {
			expected: []byte("{\"this\":\"is json\"}"),
			input:    []byte("{\"this\":\"is json\"}"),
		},
		"input is not json": {
			expected: []byte("{\"this\":\"is json\"}"),
			input: []byte(`---
this:
  is json`),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got, err := toJson(tc.input)
			if err != nil {
				t.Fatalf("There was an error converting the input to json: %s", err)
			}

			if !bytes.Equal(tc.expected, got) {
				t.Fatalf("Did not get expected result, got: %s, want: %s", string(got), string(tc.expected))
			}
		})
	}
}

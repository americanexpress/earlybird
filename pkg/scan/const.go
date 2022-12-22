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

package scan

const (
	ruleSuffix        string  = ".json"
	entropyThreshold  float64 = 4.7
	compressRegex     string  = ".(war|jar|zip|ear)$"
	convertRegex      string  = ".(docx|odt|pdf|rtf)$"
	tempRegex         string  = `(?:ebgit|ebzip|ebconv)\d+[/\\](.+$)`
	maskCharacter     string  = "*"
	overlapLength     int     = 25
	infoLevelSeverity string  = "info"
)

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

package core

import (
	cfgReader "github.com/americanexpress/earlybird/pkg/config"
)

//EarlybirdCore is the main interface for interacting with Earlybird as a package
type EarlybirdCore interface {
	ConfigInit() cfgReader.EarlybirdConfig
	StartHTTP()
	Scan()
}

//EarlybirdCfg is the global Earlybird configuration
type EarlybirdCfg struct {
	Config cfgReader.EarlybirdConfig
}

//PTRHTTPConfig is the configuration for the Earlybird REST API
type PTRHTTPConfig struct {
	HTTP       *string
	HTTPConfig *string
	HTTPS      *string
	HTTPSCert  *string
	HTTPSKey   *string
}

//PTRGitConfig is the configuration definition for Earlybird git scans
type PTRGitConfig struct {
	Repo     *string
	RepoUser *string
	Project  *string
}

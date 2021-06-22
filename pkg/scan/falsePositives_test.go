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

package scan

import "testing"

func Test_findFalsePositive(t *testing.T) {
	tests := []struct {
		name     string
		hit      Hit
		wantIsFP bool
	}{
		{
			name: "Don't skip positive finding",
			hit: Hit{
				Code:       3001,
				Filename:   "source.java",
				MatchValue: `password = "StrongFinding5813!"`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip false positive",
			hit: Hit{
				Code:       3001,
				Filename:   "source.properties",
				MatchValue: `password=None`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip xml password placeholder",
			hit: Hit{
				Code:       3007,
				Filename:   "client.axis2.xml",
				MatchValue: `<Password>password</Password>`,
			},
			wantIsFP: true,
		},
		{
			name: "Do not skip real password in xml file",
			hit: Hit{
				Code:       3007,
				Filename:   "client.axis2.xml",
				MatchValue: `<Password>abc13</Password>`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip xml passphrase placeholder",
			hit: Hit{
				Code:       3008,
				Filename:   "client.axis2.xml",
				MatchValue: `<passphrase>passphrase</passphrase>`,
			},
			wantIsFP: true,
		},
		{
			name: "Do not skip real passphrase in xml file",
			hit: Hit{
				Code:       3008,
				Filename:   "client.axis2.xml",
				MatchValue: `<passphrase>abc13</passphrase>`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip the github Url",
			hit: Hit{
				Code:       3066,
				Filename:   "redoc.standalone.js",
				MatchValue: `github.com/Mermade/oas-kit?sponsor=1","scripts":{"test":"mocha"},"browserify":{"transform":[["babelify",{"presets":["es2015"]}]]},"repository":{"url":"https://github.com/Mermade/oas-kit.git","type":"git"}`,
			},
			wantIsFP: false,
		},
		{
			name: "skip allowed encrypted CC number",
			hit: Hit{
				Code:       6008,
				Filename:   "aspire.json",
				MatchValue: `      "account_token": "EPQY2I6OK38NBKM",`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip typescript variable assignment",
			hit: Hit{
				Code:       3006,
				Filename:   "epaas-ssl.ts",
				MatchValue: `      passphrase: pass,`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip PropType definitions",
			hit: Hit{
				Code:       6008,
				Filename:   "TrackerContainer.jsx",
				MatchValue: "    account_token: PropTypes.string.isRequired,",
			},
			wantIsFP: true,
		},
		{
			name: "Skip JSON same key same value",
			hit: Hit{
				Code:     3026,
				Filename: "test.json",
				LineValue: `		"secret" : "secret"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip JS same key same value",
			hit: Hit{
				Code:     3026,
				Filename: "test.json",
				LineValue: `		'secret' : 'secret'`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip actual secret",
			hit: Hit{
				Code:     3026,
				Filename: "test.json",
				LineValue: `		"secret" : "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual client-secret",
			hit: Hit{
				Code:     3033,
				Filename: "test.yaml",
				LineValue: `		client-secret: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual clientSecret",
			hit: Hit{
				Code:     3033,
				Filename: "test.json",
				LineValue: `		clientSecret: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual 'clientSecret'",
			hit: Hit{
				Code:     3034,
				Filename: "test.json",
				LineValue: `		"clientSecret": "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual client-secret",
			hit: Hit{
				Code:     3034,
				Filename: "test.json",
				LineValue: `		"client-secret": "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip Grafana client secret",
			hit: Hit{
				Code:      3034,
				Filename:  "test.yaml",
				LineValue: `client_secret: grafana-client-secret`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip some client secret",
			hit: Hit{
				Code:      3034,
				Filename:  "test.yaml",
				LineValue: `client_secret: some-client-secret`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip actual auth-key",
			hit: Hit{
				Code:     3035,
				Filename: "test.yaml",
				LineValue: `		auth-key: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual auth-key",
			hit: Hit{
				Code:     3035,
				Filename: "test.json",
				LineValue: `		"auth-key": "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual secret-key",
			hit: Hit{
				Code:     3036,
				Filename: "test.json",
				LineValue: `		secret-key: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual secret-key",
			hit: Hit{
				Code:     3036,
				Filename: "test.json",
				LineValue: `		"secret-key": "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip credenitals placeholder",
			hit: Hit{
				Code:     3075,
				Filename: "test.yaml",
				MatchValue: `		credentials: ******`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip actual credentials in yaml",
			hit: Hit{
				Code:     3075,
				Filename: "test.yaml",
				MatchValue: `		credentials: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual credential in json",
			hit: Hit{
				Code:     3075,
				Filename: "test.json",
				MatchValue: `		"credential": "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual credential in yaml",
			hit: Hit{
				Code:     3075,
				Filename: "test.yaml",
				MatchValue: `		credential: do_not_use_this_in_production`,
			},
			wantIsFP: false,
		},
		{
			name: "skip credential with value X ",
			hit: Hit{
				Code:       3075,
				Filename:   "source.yaml",
				MatchValue: `credential: xxxxxxx`,
			},
			wantIsFP: true,
		},
		{
			name: "skip credentials with value X ",
			hit: Hit{
				Code:       3075,
				Filename:   "source.yaml",
				MatchValue: `credentials= xxxxxxx`,
			},
			wantIsFP: true,
		},
		{
			name: "skip string credential with value X ",
			hit: Hit{
				Code:       3075,
				Filename:   "source.yaml",
				MatchValue: `credential: 'xxxxxx'`,
			},
			wantIsFP: true,
		},
		{
			name: "skip string credentials with value X ",
			hit: Hit{
				Code:       3075,
				Filename:   "source.yaml",
				MatchValue: `credentials: 'xxxxxxx'`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip false positive predefined credentials (include) from Fetch API",
			hit: Hit{
				Code:     3075,
				Filename: "source.js",
				MatchValue: `		credentials=include`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip predefined credentials [omit] from Fetch API",
			hit: Hit{
				Code:     3075,
				Filename: "test.ts",
				LineValue: `		"credentials" : "omit"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip predefined credentials (same-origin) from Fetch API",
			hit: Hit{
				Code:     3075,
				Filename: "test.ts",
				LineValue: `		'credentials' : 'same-origin'`,
			},
			wantIsFP: true,
		},
		{
			name: "Donot Skip predefined credentials (same-origin) from Fetch API in java file",
			hit: Hit{
				Code:     3075,
				Filename: "test.java",
				LineValue: `		'credentials' : 'same-origin'`,
			},
			wantIsFP: false,
		},
		{
			name: "skip account token that references function call",
			hit: Hit{
				Code:       6008,
				Filename:   "check_eligibility.py",
				MatchValue: `account_token = selected_card_account_token(accounts, selected_card)`,
			},
			wantIsFP: true,
		},
		{
			name: "should not skip actual account tokens",
			hit: Hit{
				Code:       6008,
				Filename:   "check_eligibility.yaml",
				MatchValue: `"linked_account_token": "97XKNWNOVJR3U0Z",`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual secret",
			hit: Hit{
				Code:     3001,
				Filename: "test.json",
				LineValue: `		"IV_PSK_INPUT_LABEL_GENERIC_PASSWORD": "xxxxxx"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual secret",
			hit: Hit{
				Code:     3001,
				Filename: "test.properties",
				LineValue: `		IV_PSK_INPUT_LABEL_GENERIC_PASSWORD= xxxxxxxx`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual secret",
			hit: Hit{
				Code:     3001,
				Filename: "test.properties",
				LineValue: `		IV_PSK_INPUT_LABEL_GENERIC_PASSWORD= "xxxxxxxx"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual .swift secret",
			hit: Hit{
				Code:     3057,
				Filename: "test.swift",
				LineValue: `		let password = "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Don't skip actual .swift secret",
			hit: Hit{
				Code:     3057,
				Filename: "test.swift",
				LineValue: `		let password: String! = "do_not_use_this_in_production"`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip .swift FP",
			hit: Hit{
				Code:     3057,
				Filename: "test.swift",
				MatchValue: `		var txtFieldPassword: UITextField! = nil`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip swedish for password(Lösenord) as value",
			hit: Hit{
				Code:     3001,
				Filename: "source.json",
				LineValue: `		"Password": "Lösenord"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip swedish for password(Lösenord) as value in properties file",
			hit: Hit{
				Code:     3001,
				Filename: "source.properties",
				LineValue: `		"Password"= "Lösenord"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip secret: 'None'",
			hit: Hit{
				Code:       3033,
				Filename:   "source.properties",
				MatchValue: `client_secret: 'None'`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip secretkey as prompt key",
			hit: Hit{
				Code:      3036,
				Filename:  "source.txt",
				LineValue: `client_secret_key : 'e2.apigee.client.secret.key'`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip valid secret Key",
			hit: Hit{
				Code:      3036,
				Filename:  "source.txt",
				LineValue: `client_secret_key : 'valid_password'`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip valid password false positive",
			hit: Hit{
				Code:      3057,
				Filename:  "user.properties_postgress",
				LineValue: `db_password=password`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip valid password false positive with @ symbol",
			hit: Hit{
				Code:      3057,
				Filename:  "XquerySource.JSON",
				LineValue: `"connectionString": "CloudConnector;USE_ZORBA2;SourceName=eloqua;tokenType=userPassword;user=@user;password=@password"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip credentials false positive",
			hit: Hit{
				Code:      3075,
				Filename:  "imageloader.js",
				LineValue: `  USE_CREDENTIALS: 'use-credentials'`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password fields that refer to files",
			hit: Hit{
				Code:      3057,
				Filename:  "jsclassinfo.properties",
				LineValue: `mstrmojo.WH.Password: mojo/js/source/WH/Password.js`,
			},
			wantIsFP: true,
		},
		{
			name: "Ignore mixed case declaration for valid fp",
			hit: Hit{
				Code:      3057,
				Filename:  "jsclassinfo.properties",
				LineValue: "password=PASSWORD",
			},
			wantIsFP: true,
		},
		{
			name: "Skip credentials key with spaces",
			hit: Hit{
				Code:      3075,
				Filename:  "source.md",
				LineValue: `client_credentials="not a valid password"`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip credentials key with dots",
			hit: Hit{
				Code:      3075,
				Filename:  "source.md",
				LineValue: `client_credentials="key.to.password.value"`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip credentials key with valid password",
			hit: Hit{
				Code:      3075,
				Filename:  "source.md",
				LineValue: `client_credentials="secretPassword123"`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip calls to file function",
			hit: Hit{
				Code:      3075,
				Filename:  "main.tf",
				LineValue: `credentials = file(var.credPath)`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't skip calls to file when it's not a function call",
			hit: Hit{
				Code:      3075,
				Filename:  "main.tf",
				LineValue: `credentials = fileSecret1234`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip private Key as string",
			hit: Hit{
				Code:     3023,
				Filename: "source.java",
				LineValue: `		"'{{' .Data.data.privateKey | regexReplaceAll \"-----BEGIN RSA PRIVATE KEY-----\"`,
			},
			wantIsFP: true,
		},
		{
			name: "Don't Skip real private Key",
			hit: Hit{
				Code:     3023,
				Filename: "source.java",
				MatchValue: `		-----BEGIN RSA PRIVATE KEY-----`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip password as variable definitions",
			hit: Hit{
				Code:      3057,
				Filename:  "custom-environment.yaml",
				LineValue: `password: MONGO_DB_PASSWORD`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password as variable definitions in properties file",
			hit: Hit{
				Code:      3057,
				Filename:  "custom-environment.properties",
				LineValue: `password= COUCHBASE_DB_PASSWORD`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password as variable definitions",
			hit: Hit{
				Code:      3001,
				Filename:  "source.js",
				LineValue: `EPASS_ENV_VAR_SSLPASSWORD = 'KEYFILE_PASSWORD'`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password using getString",
			hit: Hit{
				Code:     3057,
				Filename: "source.scala",
				LineValue: `		val GPpassword = appConfig.getString("greenplum.jdbc.pwd")`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password that is a function",
			hit: Hit{
				Code:     3057,
				Filename: "source.xml",
				LineValue: `		String password = context.decrypt("%%ARCHIVAL_DB_PWD%%");`,
			},
			wantIsFP: true,
		},
		{
			name: "Skip password that is a function",
			hit: Hit{
				Code:     3057,
				Filename: "source.xml",
				LineValue: `		password = generatePolicyPassword(context, identity, "AXP Directory")`,
			},
			wantIsFP: true,
		},
		{
			name: "Do not skip password in logger",
			hit: Hit{
				Code:     3057,
				Filename: "source.xml",
				LineValue: `		logger.debug("Password Rule Library : generateNotAPolicyPassword :hasDigit :" + hasDigit);`,
			},
			wantIsFP: false,
		},
		{
			name: "Skip auth Key as value at the end",
			hit: Hit{
				Code:      3035,
				Filename:  "source.js",
				LineValue: `authKey: GITHUB_AUTH_KEY`,
			},
			wantIsFP: true,
		},
		{
			name: "Do not skip actual auth Key",
			hit: Hit{
				Code:      3035,
				Filename:  "source.js",
				LineValue: `authKey: V3RyStr0nG@uThk3Y`,
			},
			wantIsFP: false,
		},
		{
			name: "Donot skip password in js file",
			hit: Hit{
				Code:      3057,
				Filename:  "source.js",
				LineValue: `Password=my_logo_i_passwd>what`,
			},
			wantIsFP: false,
		},
		{
			name: "Donot skip credentials with just 4 digits ",
			hit: Hit{
				Code:       3075,
				Filename:   "source.txt",
				MatchValue: `CREDENTIALS=2727`,
			},
			wantIsFP: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsFP := findFalsePositive(tt.hit); gotIsFP != tt.wantIsFP {
				t.Errorf("findFalsePositive() = %v, want %v", gotIsFP, tt.wantIsFP)
			}
		})
	}
}

func Test_falsePositiveInclusivity(t *testing.T) {
	tests := []struct {
		name     string
		hit      Hit
		wantIsFP bool
	}{
		{
			name: "Ignore references to ticketmaster",
			hit: Hit{
				Code:       20014,
				Filename:   "file.properties",
				MatchValue: "ticketmaster",
			},
			wantIsFP: true,
		},
		{
			name: "Ignore references to mastercard",
			hit: Hit{
				Code:       20014,
				Filename:   "file.properties",
				MatchValue: "mastercard",
			},
			wantIsFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsFP := findFalsePositive(tt.hit); gotIsFP != tt.wantIsFP {
				t.Errorf("findFalsePositive() = %v, want %v", gotIsFP, tt.wantIsFP)
			}
		})
	}
}

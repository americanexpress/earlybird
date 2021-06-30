package postprocess

import "testing"

type params struct {
	matchValue string
	lineValue  string
}

var testsSkipSameKeyValuePasswords []struct {
	name       string
	params     params
	wantIgnore bool
} = []struct {
	name       string
	params     params
	wantIgnore bool
}{
	{
		name: "Skip same key/value passwords as case insensitive in properties file",
		params: params{
			matchValue: "PASSWORD_DB : password_db",
			lineValue:  "PASSWORD_DB : password_db",
		},
		wantIgnore: true,
	},
	{
		name: "Skip same key/value passwords as case insensitive in yaml file in lineValue",
		params: params{
			matchValue: "PASSWORD : db_password",
			lineValue:  "DB_PASSWORD : db_password",
		},
		wantIgnore: true,
	},
	{
		name: "Skip same key/value passwords that match alphanumerically in lineValue",
		params: params{
			matchValue: "PASSWORD : $db.password",
			lineValue:  "DB_PASSWORD : $db.password",
		},
		wantIgnore: true,
	},
	{
		name: "Skip same key/value secret in properties file",
		params: params{
			matchValue: "PASSWORD : $db.password",
			lineValue:  "DB_PASSWORD : $db.password",
		},
		wantIgnore: true,
	},
	{
		name: "Skip same key/value secret in yaml file in lineValue",
		params: params{
			matchValue: "SECRET: api.Secret",
			lineValue:  "APISECRET: api.Secret",
		},
		wantIgnore: true,
	},
	{
		name: "Skip same key/value secret in json file",
		params: params{
			matchValue: "\"SECRET\": \"SECRET\"",
			lineValue:  "\"SECRET\": \"SECRET\"",
		},
		wantIgnore: true,
	},
	{
		name: "Do not skip, real password finding",
		params: params{
			matchValue: "password_couchbase: VeryStrong857#",
			lineValue:  "password_couchbase: VeryStrong857#",
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real password finding",
		params: params{
			matchValue: "password: VeryStrong857#",
			lineValue:  "couchbase_password: VeryStrong857#",
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real secret finding",
		params: params{
			matchValue: "secret: VeryStrong857#",
			lineValue:  "secret: VeryStrong857#",
		},
		wantIgnore: false,
	},
	{
		name: "Skip same key/value when match value and line value are different",
		params: params{
			matchValue: "Secret=apiAppSecret",
			lineValue:  "api.appSecret=apiAppSecret",
		},
		wantIgnore: true,
	},
}

func TestSkipPasswordSameKeyValue(t *testing.T) {
	for _, tt := range testsSkipSameKeyValuePasswords {
		t.Run(tt.name, func(t *testing.T) {
			gotIgnore := SkipSameKeyValue(tt.params.matchValue, tt.params.lineValue)
			if gotIgnore != tt.wantIgnore {
				t.Errorf("SkipSameKeyValue() gotIgnore = %v, want %v", gotIgnore, tt.wantIgnore)
			}
		})
	}
}

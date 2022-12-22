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

package git

import (
	"io"
	"strings"
	"testing"
)

func Test_splitDiffs(t *testing.T) {
	gitReader := strings.NewReader(`commit 719709695ab1041c8cde51b721cdc4e63cbac389 (HEAD -> example, origin/example)
		Author: Foo Bar <foo@bar.com>
		Date:   Fri Jan 17 11:34:35 2020 -0700
		
			feature(test): Important GIT diff test
			
			Test parsing sample git log
		
		diff --git a/sample.go b/sample.go
		index 3c3108d..f53837c 100644
		--- a/sample.go
		+++ b/sample.go
		@@ -1,3+1,3 @@ FIRST LINE
		-       DELETED THIS LINE
		+       ADDED THIS LINE
		+       ADDED THIS LINE AGAIN
	`)
	diff := Diff{}
	type args struct {
		r io.Reader
		l List
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Split diffs in buffer",
			args: args{
				r: gitReader,
				l: &diff,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := splitDiffs(tt.args.r, tt.args.l); (err != nil) != tt.wantErr {
				t.Errorf("splitDiffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

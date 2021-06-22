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

package writers

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//WriteJSON Outputs an object as a JSON blob to an output file or console
func WriteJSON(v interface{}, fileName string) (s string, err error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return "", err
	}
	if fileName == "" {
		_, err = os.Stdout.Write(b)
	} else {
		err = ioutil.WriteFile(fileName, b, 0666)
	}
	return string(b), err

}

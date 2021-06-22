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
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"

	"github.com/americanexpress/earlybird/pkg/scan"
)

//WriteCSV outputs Earlybird hit findings from the worker channel to files or console
func WriteCSV(hits <-chan scan.Hit, fileName string) (err error) {
	var printedheader bool
	// If no filename was passed in, just print to stdout
	if fileName == "" {
		for hit := range hits {
			printedheader, err = hitToCSVWriter(printedheader, hit, os.Stdout)
			if err != nil {
				return err
			}
		}
		return nil
	}

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	// Close the file after the function ends
	defer csvFile.Close()
	for hit := range hits {
		printedheader, err = hitToCSVWriter(printedheader, hit, csvFile)
		if err != nil {
			return err
		}
	}
	//Get actual file stats for file size
	fi, err := csvFile.Stat()
	if err != nil {
		return err
	}
	fmt.Println(fi.Size(), " bytes written to ", fileName)
	return nil
}

func hitToCSVWriter(printedheader bool, hit scan.Hit, output io.Writer) (bool, error) {
	if printedheader { //Check if we already printed our file header
		err := gocsv.MarshalWithoutHeaders(&[]scan.Hit{hit}, output) // Use this to save the CSV back to the file
		if err != nil {
			return printedheader, err
		}
	} else { //add file header
		err := gocsv.Marshal(&[]scan.Hit{hit}, output) // Get all clients as CSV string
		if err != nil {
			return false, err
		}
		return true, err
	}
	return printedheader, nil
}

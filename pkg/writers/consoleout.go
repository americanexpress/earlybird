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
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/americanexpress/earlybird/pkg/scan"
)

type issue struct {
	k string
	c int
}

var issues = make(map[string]int)

//WriteConsole streams hits from the result channel to the command line or target file
func WriteConsole(hits <-chan scan.Hit, fileName string, showFullLine bool) error {
	// If no filename was passed in, just print to stdout
	if fileName == "" {
		// Store a record of the number of each threat found
		i := 1
		for hit := range hits {
			fmt.Println(hitToConsole(hit, i, showFullLine))
			issues[hit.Caption]++
			i++
		}
	} else {
		err := hitsToFile(hits, fileName, showFullLine)
		if err != nil {
			log.Fatal("Failed to write results to file", err)
		}
	}
	displayIssues()
	return nil
}

func hitsToFile(hits <-chan scan.Hit, fileName string, showFullLine bool) error {
	// Store a record of the number of each threat found
	i := 1
	f, createErr := os.Create(fileName)
	if createErr != nil {
		return createErr
	}
	writer := bufio.NewWriter(f)

	// Close the file after the function ends
	defer func() {
		writer.Flush()
		f.Sync()
		f.Close()
	}()

	for hit := range hits {
		_, err := writer.WriteString(hitToConsole(hit, i, showFullLine))
		if err != nil {
			return err
		}

		issues[hit.Caption]++
		i++
	}

	//Get actual file stats for file size
	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(fi.Size(), outputBytesWritten, fileName)
	return nil
}

func displayIssues() {
	//Sort out values
	keyvals := make([]issue, 0, len(issues))
	for k, v := range issues {
		keyvals = append(keyvals, issue{k, v})
	}
	sort.Slice(keyvals, func(i, j int) bool {
		a, b := keyvals[i], keyvals[j]
		// we want to sort by count DESCENDING, but
		// alphabetically in the case of a tie.
		// sort.Slice takes a function "less" that assumes you're sorting in
		// ASCENDING order, so we need to greater than for the counts
		// to reverse the order
		return a.c > b.c || a.c == b.c && a.k < b.k
	})

	//Print out our findings summary
	fmt.Println(outputTotalIssuesFnd)
	var total int
	for _, issue := range keyvals {
		fmt.Printf("\t%5d %s\n", issue.c, issue.k)
		total += issue.c
	}
	fmt.Printf(outputTotalIssues, total)
}

func hitToConsole(hit scan.Hit, progress int, showFullLine bool) string {
	var sb strings.Builder
	sb.WriteString(columnFinding + " " + strconv.Itoa(progress) + ":")
	sb.WriteString(outputIndent + columnCode + ": " + strconv.Itoa(hit.Code))
	sb.WriteString(outputIndent + columnFileName + ": " + hit.Filename)
	sb.WriteString(outputIndent + columnCaption + ": " + hit.Caption)
	sb.WriteString(outputIndent + columnCategory + ": " + hit.Category)
	sb.WriteString(outputIndent + columnLine + ": " + strconv.Itoa(hit.Line))
	sb.WriteString(outputIndent + columnValue + ": " + printableASCII(hit.MatchValue))
	if showFullLine {
		sb.WriteString(outputIndent + columnLineValue + ": " + printableASCII(hit.LineValue))
	}
	sb.WriteString(outputIndent + columnSeverity + ": " + hit.Severity)
	sb.WriteString(outputIndent + columnConfidence + ": " + hit.Confidence)
	sb.WriteString(outputIndent + columnLabels + ": " + displayCWE(hit.Labels))
	sb.WriteString(outputIndent + columnCWE + ": " + displayCWE(hit.CWE))
	if hit.Solution != "" {
		sb.WriteString(outputIndent + columnSolution + ": " + hit.Solution)
	}
	sb.WriteString("\n")
	return sb.String()
}

func displayCWE(cwe []string) (result string) {
	if len(cwe) > 0 {
		result = strings.Join(cwe, outputArraySeparator)
	} else {
		result = outputNone
	}
	return result
}

// strip control characters, DELETE, and non-ASCII unicode from the string.
func printableASCII(s string) string {
	printable := make([]byte, 0, len(s))
	for _, b := range []byte(s) {
		if b >= 32 && b < 127 {
			printable = append(printable, b)
		}
	}
	return string(printable)
}

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

package file

import (
	"archive/zip"
	"bufio"
	"bytes"
	"code.sajari.com/docconv"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/scan"
	"github.com/americanexpress/earlybird/pkg/utils"
	"github.com/americanexpress/earlybird/pkg/wildcard"
)

var (
	ignoreFiles    = [...]string{".ge_ignore"}
	ignorePatterns []string
)

//MultipartToScanFiles converts the multipart file upload into Earlybird files
func MultipartToScanFiles(files []*multipart.FileHeader, cfg cfgreader.EarlybirdConfig) (fileList []scan.File, err error) {
	ignorePatterns = getIgnorePatterns(cfg.SearchDir, cfg.IgnoreFile, cfg.VerboseEnabled)

	for _, fheader := range files {
		myfile, err := fheader.Open()
		if err != nil {
			return fileList, err
		}
		// Per the HTTP spec, The filename directive of multipart form data will have it's path information stripped https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition.
		// .ge_ignore file only works on absolute paths, not the basename of a file
		// client will send filepath as base64 encoded and earlybird will decode to get the full path
		// a `/` in order to have the .ge_ignore rules apply for files scanned by HTTP request
		fileNameBytes, err := base64.StdEncoding.DecodeString(fheader.Filename)
		var fileName string
		// If filename is passed as utf-8 string then base64 decode will throw error
		if err != nil {
			fileName = fheader.Filename // Support utf-8 filename as backward compatibility
		} else {
			fileName = string(fileNameBytes) // Use base64 decoded value
		}
		fileNameWithPathPrefix := fileName
		pathSeparator := "/"
		if !strings.HasPrefix(fileName, pathSeparator) {
			fileNameWithPathPrefix = pathSeparator + fileNameWithPathPrefix
		}
		//Skip file with extensions Earlybird ignores
		if isIgnoredFile(fileNameWithPathPrefix, cfg.SearchDir) {
			continue
		}

		//Start of file upload parsing (indepth comments in scanUtil.go)
		curFile := scan.File{
			Name: fileNameWithPathPrefix,
			Path: "buffer",
		}

		var line scan.Line
		reader := bufio.NewReader(myfile)
		for {
			var buffer bytes.Buffer

			var l []byte
			var isPrefix bool
			for {
				l, isPrefix, err = reader.ReadLine()
				buffer.Write(l)
				//Reached the end of the line, stop reading.
				if !isPrefix {
					break
				}
				// EOF, break
				if err != nil {
					break
				}
			}
			if err == io.EOF {
				break
			}
			lineText := buffer.String()
			line.LineNum = line.LineNum + 1
			line.LineValue = lineText
			line.FilePath = curFile.Path
			line.FileName = fileNameWithPathPrefix
			curFile.Lines = append(curFile.Lines, line)
		}
		if err != io.EOF {
			return fileList, err
		}

		myfile.Close()
		fileList = append(fileList, curFile)
	}
	return
}

//GetGitFiles Builds the list of staged or tracked files
func GetGitFiles(fileType string, cfg *cfgreader.EarlybirdConfig) (fileContext Context, err error) {
	ignorePatterns = getIgnorePatterns(cfg.SearchDir, cfg.IgnoreFile, cfg.VerboseEnabled)

	var (
		output       []byte
		compressList []scan.File
		convertList  []scan.File
		fileList     []scan.File
	)

	if fileType == utils.Tracked {
		output, err = exec.Command("git", "ls-tree", "--full-tree", "-r", "--name-only", "HEAD").Output()
	} else if fileType == utils.Staged {
		output, err = exec.Command("git", "--no-pager", "diff", "--name-only", "--staged").Output()
	}

	if err != nil {
		log.Println(notTrackedDir)
		return fileContext, err
	}

	fileList = parseGitFiles(output, cfg.VerboseEnabled, cfg.MaxFileSize)
	compressList, fileList = separateCompressedAndUncompressed(fileList)
	compressList, fileContext.CompressPaths, err = GetCompressedFiles(compressList, cfg.SearchDir) //Get the files within our compressed list
	if err != nil {
		return fileContext, err
	}
	fileContext.Files = append(fileList, compressList...)
	convertList, fileContext.ConvertPaths = GetConvertedFiles(fileContext.Files) //Get the files that need to be converted and convert them to plaintext
	fileContext.Files = append(fileContext.Files, convertList...)
	fileContext.IgnorePatterns = ignorePatterns
	return fileContext, nil
}

func parseGitFiles(out []byte, verbose bool, maxFileSize int64) (fileList []scan.File) {
	var curFile scan.File
	// Convert byteArray to string
	gitFiles := string(out)
	if len(gitFiles) < 1 {
		log.Println(gitErr)
	} else {
		// Parse the directory tree into a slice of scan.File objects
		scanner := bufio.NewScanner(strings.NewReader(gitFiles))
		for scanner.Scan() {
			curFile.Path = scanner.Text()
			curFile.Name = filepath.Base(scanner.Text())
			if fileExists := Exists(curFile.Path); fileExists {
				pathIsDirectory, dirErr := isDirectory(curFile.Path)
				if dirErr != nil {
					log.Println(dirErr)
				}
				if !pathIsDirectory && getFileSizeOK(curFile.Path, maxFileSize) {
					if verbose {
						log.Println("Reading file ", curFile.Path)
					}
					fileList = append(fileList, curFile)
				}
			}
		}
	}
	return fileList
}

//GetFiles Build the list of files
func GetFiles(searchDir, ignoreFile string, verbose bool, maxFileSize int64) (fileContext Context, err error) {
	ignorePatterns = getIgnorePatterns(searchDir, ignoreFile, verbose)
	fileList := make([]scan.File, 0)
	var curFile scan.File
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error reading directory: ", err)
		}
		if !isIgnoredFile(path, searchDir) {
			// Ignore the path if it's a directory
			pathIsDirectory, isDirErr := isDirectory(path)
			if !pathIsDirectory {
				if isDirErr != nil && verbose {
					log.Println("Error checking if path is directory")
				}
				if getFileSizeOK(path, maxFileSize) {
					curFile.Name = f.Name()
					curFile.Path = path
					fileList = append(fileList, curFile)
					if verbose {
						log.Println("Reading file ", curFile.Path)
					}
				} else {
					fileContext.SkippedFiles = append(fileContext.SkippedFiles, path)
					if verbose {
						log.Println("Ignoring", path, ". Filesize is too large.")
					}
				}
			}
		} else {
			fileContext.SkippedFiles = append(fileContext.SkippedFiles, path)
			if verbose {
				log.Println("Ignoring", path, ". File blacklisted.")
			}
		}
		return err
	})
	if err != nil {
		return fileContext, err
	}

	var compressList, convertList []scan.File
	compressList, fileList = separateCompressedAndUncompressed(fileList)
	compressList, fileContext.CompressPaths, err = GetCompressedFiles(compressList, searchDir) //Get the files within our compressed list
	if err != nil {
		return fileContext, err
	}
	fileContext.Files = append(fileList, compressList...)
	convertList, fileContext.ConvertPaths = GetConvertedFiles(fileContext.Files) //Get the files that need to be converted and convert them to plaintext
	fileContext.Files = append(fileContext.Files, convertList...)
	fileContext.IgnorePatterns = ignorePatterns
	return fileContext, nil
}

//GetFileFromStream Builds a file as a collection of lines from the input stream.
// This will be fed to the scan modules.
func GetFileFromStream(cfg *cfgreader.EarlybirdConfig) []scan.File {
	// Read Stdin
	scanner := bufio.NewScanner(os.Stdin)
	// The scan modules will expect a list of Files, so create that list with just one
	fileList := make([]scan.File, 0)
	curFile := scan.File{
		Name: "buffer",
		Path: "buffer",
	}

	var line scan.Line
	// Set the stage for escaping lines with the EARLYBIRD-IGNORE annotation
	ignoreAnnotationInLine := false
	nextLineIgnored := false

	// Read the initial line
	for scanner.Scan() {
		lineText := scanner.Text()
		// If we find the EARLYBIRD-IGNORE annotation...
		if scan.IsIgnoreAnnotation(cfg, lineText) {
			lineText = ""
			ignoreAnnotationInLine = true
			nextLineIgnored = true
		}
		// If we had previously found the EARLYBIRD-IGNORE annotation...
		if nextLineIgnored {
			lineText = ""
		}
		line.LineNum = line.LineNum + 1
		line.LineValue = lineText
		line.FilePath = curFile.Path
		curFile.Lines = append(curFile.Lines, line)

		// Cancel any flags to ignore the next line for the next iteration
		if nextLineIgnored && !ignoreAnnotationInLine {
			nextLineIgnored = false
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Reading standard input:", err)
	}
	fileList = append(fileList, curFile)
	return fileList
}

//GetFileSize returns the file size of target file
func GetFileSize(path string) (size int64, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

// Make sure the filesize is lower than the MAX_FILE_SIZE threshold so bufio doesn't fail
func getFileSizeOK(path string, maxFileSize int64) bool {
	size, err := GetFileSize(path)
	return err == nil && size != 0 && (size < maxFileSize || hasCompressionExtension(path))
}

func hasCompressionExtension(path string) bool {
	switch filepath.Ext(path) {
	case ".war", ".jar", ".zip", ".ear":
		return true
	default:
		return false
	}
}

// Read in .ge_ignore file and ignore files matching the patterns
func getIgnorePatterns(filePath, ignoreFile string, verbose bool) (ignorePatterns []string) {
	ignorePatterns = append(ignorePatterns, "*.git/*")

	// Loop through the files defined to contain ignore patterns (.ge_ignore, .gitignore, etc.)
	for _, ignoreFile := range ignoreFiles {
		actualFilePath := path.Join(filePath, ignoreFile)
		if Exists(actualFilePath) {
			file, err := os.Open(actualFilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var line, firstChar string
			for scanner.Scan() {
				line = scanner.Text()

				// Ignore comment lines (starting with #)
				runes := []rune(line)
				firstChar = string(runes[0:1])
				if !(firstChar == "#") && strings.Trim(line, " ") != "" {
					ignorePatterns = append(ignorePatterns, line)
				}
			}
		}
	}

	if ignoreFile != "" {
		file, err := os.Open(ignoreFile)
		if err != nil {
			log.Println("Failed to open ignore file", err)
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var line, firstChar string
			for scanner.Scan() {
				line = scanner.Text()

				// Ignore comment lines (starting with #)
				runes := []rune(line)
				firstChar = string(runes[0:1])
				if !(firstChar == "#") && strings.Trim(line, " ") != "" {
					ignorePatterns = append(ignorePatterns, line)
				}
			}
		}
	}

	if verbose {
		log.Println("Ignore pattern: ", strings.Join(ignorePatterns, ", "))
	}
	return ignorePatterns
}

// If the file matches a pattern in one of the ignore files, return true
func isIgnoredFile(fileName string, fileRoot string) bool {
	// ignore root directory when checking ignore matching
	trimmedName := strings.Replace(fileName, fileRoot, "", 1)
	for _, pattern := range ignorePatterns {
		if wildcard.PatternMatch(trimmedName, pattern) {
			return true
		}
	}
	return false
}

// Check a path to see if it's a directory
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

//GetWD Gets the current working directory
func GetWD() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir, err
}

//IsEmpty Check to see if a directory is empty
func IsEmpty(path string) (bool, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fileHandle.Close()

	_, err = fileHandle.Readdirnames(1) // Or f.Readdir(1)
	return err == io.EOF, err
}

//Exists Check to see if a path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func separateCompressedAndUncompressed(files []scan.File) (compressed, uncompressed []scan.File) {
	for _, f := range files {
		if scan.CompressPattern.MatchString(f.Path) {
			compressed = append(compressed, f)
		} else {
			uncompressed = append(uncompressed, f)
		}
	}
	return compressed, uncompressed
}

//GetCompressedFiles provides all the files contained within compressed files
func GetCompressedFiles(files []scan.File, rootPath string) (newfiles []scan.File, compresspaths []string, err error) {
	//check if file list contains compressed files, if so, scan their contents
	for _, file := range files {
		//Unpack and append to file list
		tmppath, err := ioutil.TempDir("", "ebzip")
		if err != nil {
			return newfiles, compresspaths, err
		}
		compresspaths = append(compresspaths, tmppath)
		filenames, err := Uncompress(file.Path, tmppath)
		if err != nil {
			return newfiles, compresspaths, err
		}
		for _, subfile := range filenames {
			if !isIgnoredFile(subfile, rootPath) && !scan.CompressPattern.MatchString(subfile) {
				var curFile scan.File
				curFile.Path = file.Path + "/" + filepath.Base(subfile)
				curFile.Name = subfile //Build view file name in format: file.zip/contents/file
				newfiles = append(newfiles, curFile)
			}
		}
	}
	return newfiles, compresspaths, nil
}

//Uncompress decompresses zip files safely
func Uncompress(src string, dest string) (filenames []string, err error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		//Create anonymous function to avoid leaving too many files open
		copyZipFile := func(f *zip.File) ([]string, error) {
			rc, err := f.Open()
			if err != nil {
				return filenames, err
			}

			defer rc.Close()

			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name)

			// Check for ZipSlip exploit
			if !strings.HasPrefix(fpath, filepath.Clean(dest+string(os.PathSeparator))) {
				return filenames, fmt.Errorf("%s: illegal file path", fpath)
			}

			filenames = append(filenames, fpath)

			if f.FileInfo().IsDir() {
				// Make Folder
				err = os.MkdirAll(fpath, os.ModePerm)
				if err != nil {
					return filenames, err
				}
			} else {
				if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
					return filenames, err
				}
				body, err := ioutil.ReadAll(rc)
				if err != nil {
					return filenames, err
				}
				err = ioutil.WriteFile(fpath, body, 0644)
				if err != nil {
					return filenames, err
				}
			}
			return filenames, nil
		}

		filenames, err = copyZipFile(f)
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

//GetConvertedFiles converts files into plaintext
func GetConvertedFiles(files []scan.File) (convertedFiles []scan.File, convertedPaths []string) {
	var toBeConverted []scan.File
	for _, f := range files {
		if scan.ConvertPattern.MatchString(f.Path) {
			toBeConverted = append(toBeConverted, f)
		}
	}

	for _, file := range toBeConverted {
		tmppath, err := ioutil.TempDir("", "ebconv")
		fpath := filepath.Join(tmppath, file.Name)

		// Get content from the file as a string
		content, err := docconv.ConvertPath(file.Path)
		if err != nil {
			log.Printf("Error converting %s, file not scanned\n", file.Path)
			continue
		}

		// Write content to new temp file
		err = ioutil.WriteFile(fpath, []byte(content.Body), 0644)
		if err != nil {
			log.Printf("Error writing converted file %s, file not scanned\n", file.Path)
			continue
		}

		convertedPaths = append(convertedPaths, tmppath)
		var convertedFile scan.File
		convertedFile.Path = fpath
		convertedFile.Name = file.Path
		convertedFiles = append(convertedFiles, convertedFile)
	}

	return convertedFiles, convertedPaths

}

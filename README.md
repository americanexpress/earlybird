![Logo](docs/GoEarlyBird-logo_sm.png)

EarlyBird is a sensitive data detection tool capable of scanning source code repositories for clear text password violations, PII, outdated cryptography methods, key files and more. It can be used to scan remote git repositories, local files or directories or as a pre-commit step.

## Installation
### Linux & Mac
Running the `build.sh` script will produce a binary for each OS, while the `install.sh` script will install Earlybird on your system. This will create a `.go-earlybird` directory in your home directory with all the configuration files. Finally installing `go-earlybird` as an executable in `/usr/local/bin/`.
```
./build.sh && ./install.sh
```

### Windows
Running `build.bat` will produce your binaries while the `install.bat` script will create a 'go-earlybird' directory in `C:\Users\[my user]\App Data\`, and copy the required configurations there.  This script will also install `go-earlybird.exe` as an executable in the App Data directory (which should be in your path).

```
build.bat && install.bat
```

## Usage
To launch a basic EarlyBird scan against a directory:
```
$ go-earlybird --path=/path/to/directory
```
```
$ go-earlybird.exe --path=C:\path\to\directory
```
or to scan a remote git repo:
```
$ go-earlybird --git=https://github.com/americanexpress/earlybird
```
[Click here for Detailed Usage instructions.](./docs/USAGE.md)


## Documentation
 - [Usage - How do I use Earlybird?](./docs/USAGE.md)
 - [Modules - What is a Module? How do I create one?](./docs/MODULES.md)
 - [Hooks - How do I use Earlybird as Pre-Commit Hook?](./docs/HOOKS.md)
 - [REST API - How do I use Earybird as REST API?](./docs/REST.md)
 - [False Positives - How are they managed? How do I filter them?](./docs/FALSEPOSITIVES.md)
 - [Labels - What are labels? How do I create my own?](./docs/LABELS.md)
 - [Ignore - How do I skip lines or files intentionally?](./docs/IGNORE.md)


## Why Are We Doing This?
The MITRE Corporation provides a catalog of [Common Weakness Enumerations](https://cwe.mitre.org/index.html) (CWE), documenting issues that should be avoided.  Some of the relevant CWEs that are handled by the use of EarlyBird include:
 - [CWE-798 - Use of Hardcoded Credentials](https://cwe.mitre.org/data/definitions/798.html)
 - [CWE-259 - Use of Hardcoded Password](https://cwe.mitre.org/data/definitions/259.html)
 - [CWE-321 - Use of Hardcoded Cryptographic Key](https://cwe.mitre.org/data/definitions/321.html)
 - [CWE-257 - Storing Password in a Recoverable Format](https://cwe.mitre.org/data/definitions/257.html)
 - [CWE-312 - Cleartext Storage of Sensitive Information](https://cwe.mitre.org/data/definitions/312.html)
 - [CWE-327 - Use of Broken or Risky Cryptographic Algorithm](https://cwe.mitre.org/data/definitions/327.html)
 - [CWE-338 - Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)](https://cwe.mitre.org/data/definitions/338.html)
 - [CWE-615 - Information Exposure Through Comments](https://cwe.mitre.org/data/definitions/615.html)
 - [CWE-546 - Suspicious Comments](https://cwe.mitre.org/data/definitions/546.html)
 - [CWE-521 - Weak Password Requirements](https://cwe.mitre.org/data/definitions/521.html)

---

## Contributing
We welcome your interest in the American Express Open Source Community on Github. Any Contributor to
any Open Source Project managed by the American Express Open Source Community must accept and sign
an Agreement indicating agreement to the terms below. Except for the rights granted in this 
Agreement to American Express and to recipients of software distributed by American Express, You
reserve all right, title, and interest, if any, in and to your contributions. Please
[fill out the Agreement](https://cla-assistant.io/americanexpress/earlybird).

## Usage
The executable can be called from the command line with the following syntax:
```
~/go/src/gearlybird (master ✘)✭ ᐅ go-earlybird --help
Usage of go-earlybird:
  -config string
    	Directory where configuration files are stored (default "/Users/janedoe/.go-earlybird/")
  -display-confidence string
    	Lowest confidence level to display [ critical | high | medium | low ] (default "high")
  -display-severity string
    	Lowest severity level to display [ critical | high | medium | low ] (default "medium")
  -enable value
    	Enable individual scanning modules [ ccnumber | content | filename | password-secret ]
  -fail-confidence string
    	Lowest confidence level at which to fail [ critical | high | medium | low ] (default "high")
  -fail-severity string
    	Lowest severity level at which to fail [ critical | high | medium | low ] (default "high")
  -file string
    	Output file -- e.g., 'go-earlybird --file=/home/jdoe/myfile.csv'
  -format string
    	Output format [ console | json | csv ] (default "console")
  -git string
    	Full URL to a git repo to scan e.g. github.com/user/repo
  -git-commit-stream
    	Use stream IO of Git commit log as input instead of file(s) -- e.g., 'cat secrets.text > go-earlybird'
  -git-project string
    	Full URL to a github organization or bitbucket project to scan e.g. github.com/org
  -git-staged
    	Scan only git staged files
  -git-tracked
    	Scan only git tracked files
  -git-user string
    	If the git repository is private, enter an authorized username
  -http string
    	Listen IP and Port for HTTP API e.g. 127.0.0.1:8080
  -http-config string
    	Path to webserver config JSON file
  -https string
    	Listen IP and Port for HTTPS/2 API e.g. 127.0.0.1:8080 (Don't forget the https-cert and https-key flags)
  -https-cert string
    	Certificate file for TLS
  -https-key string
    	Private key file for TLS
  -ignore-fp-rules
    	Ignore the false positive post-process rules
  -ignorefile string
    	Patterns File (including wildcards) for files to ignore.  (e.g. *.jpg) (default "/Users/jhans12/.ge_ignore")
  -max-file-size int
    	Maximum file size to scan (in bytes) (default 10240000)
  -path string
    	Directory to scan (defaults to CWD) -- ABSOLUTE PATH ONLY (default "/Users/jhans12/go/src/gearlybird")
  -show-full-line
    	Display the full line where the pattern match was found (warning: this can be dangerous with minified script files)
  -show-rules-only
    	Display rules that would be run, but do not execute a scan
  -skip-comments
    	Skip scanning comments in files -- applies only to the 'content' module
  -stream
    	Use stream IO as input instead of file(s)
  -suppress
    	Suppress reporting of the secret found (important if output is going to Slack or other logs)
  -update
    	Update module configurations
  -verbose
    	Reports details about file reads
  -workers int
    	Set number of workers. (default 100)
  -worksize int
    	Set Line Wrap Length. (default 2500)
  -module-config-file string
        Absolute path to a json or yaml file for per module level config -- {"modules": { "aModule": { "display_severity": "medium" } } }
  ```

### Creating New Modules:
New modules can be added via json rules files in the user's go-earlybird configuration directory.  Simply add a new json file into this directory (e.g. `custom-rules.json`) with the following structure, and Go-EarlyBird will detect and load the rules.  Keeping these custom rules in a separate file will ensure they do not get overwritten when Go-EarlyBird is updated.
```
{
    "Searcharea": "<Where in the file should the scan search?  Supports 'body' or 'filename'>",
    "rules": [
    {
      "Code": 1,
      "Pattern": "<Regexp pattern>",
      "Caption": "<A description of the finding (e.g., password, PII value, etc.)>",
      "Solution": "<Reference ID from solutions.json",
      "Category": "<The type of finding>",
      "Severity": <Value of 1-4 with 1 being critical and 4 being low>,
      "Confidence": <Value of 1-4 with 1 being critical and 4 being low>,
      "Postprocess": "<functions that perform extra validation on a hit value: password, mod10, entropy, ssn>",
      "Example": "password='xxx'",
      "CWE": ["CWE-XXX"]
    }
  ]
}
```
We recommend using a unique, integer-only approach to defining the `Code` field.

Custom labels and custom solutions can be added in the same manner.


---
&nbsp;


## False Positive Detection

EarlyBird has a rules engine for excluding false positives from the results.  Each rule in `false-positives.yaml` is tied to one or more scan rules (using the `Codes` field).

The `Pattern` field is a regular expression that is evaluated against any hit that matches the `Code`, as long as the file containing that hit has an extension matching a value in the `FileExtensions` value (if that value is empty, all file extensions will be considered).

```
---
rules:
- Codes:
  - 3013
  Pattern: "(000-000-0000)"
  FileExtensions: []
  Description: Ignore a false positive phone number
- Codes:
  - 4005
  - 3022
  Pattern: ".*"
  FileExtensions:
  - ".md"
  - ".txt"
  - ".doc"
  - ".pdf"
  - ".docx"
  - ".csv"
  - ".html"
  - ".htm"
  Description: Ignore deprecated crypto in documents
```

In the examples above:
1. Any hit found with rule 3013 (looking for 10-digit phone number patterns), that matches all zeroes (000-000-0000), in any type of file will be ignored
2. Any hit found with rules 4005 or 3022 (looking for indicators of deprecated crypto method use like 3DES or JUICE) in document files will be ignored

---
&nbsp;

## Labels

By modifying `.go-earlybird/labels.yaml`, or adding a separate file, labels can be added to the hits based on context.  This can be done either on a same line search or on a scan of context throughout the whole file.  For example, the following line will add the `oracle` label to the returned hit:



## <a name="running"></a> Running Go-EarlyBird

### <a name="standalone"></a> Standalone
Assuming the setup script was run, you can kick off the application by running 'go-earlybird' / 'go-earlybird.exe' / 'go-earlybird-linux (mac / windows / linux).  See the *Usage* section below

If Go is installed, the project can be downloaded and run with `go run go-earlybird.go`

### Streamed / Piped input
Using the `-stream` flag, users can stream or pipe file contents to 'go-earlybird'.  

```
ᐅ go-earlybird -stream < /path/to/file
```
... or:
```
ᐅ cat /path/to/file | go-earlybird -stream
```

### HTTP API
```
ᐅ go-earlybird --http 0.0.0.0:3000
```
`/scan` will accept a multi-part upload and scan the contents, returning json output.

The normal HTTP listener will operate on HTTP/1.1.  Go-EarlyBird can be run as HTTPS/2 with the `-https [ip:port]` flag.  Note that this also requires the `-https-cert [/path/to/cert]` and `-https-key [/path/to/key]` parameters.

The simple webserver configuration file can be found in the local config directory (`~/.go-earlybird/webserver.json` or `C:\Users\[me]\AppData\go-earlybird\webserver.json`).  A separate config file can be specified using the `-http-config [/path/to/configfile]` flag.


### Local Git Scanning
With the flag `-git-staged` or `-git-tracked`, Go-EarlyBird can limit its scan to only look at files that are staged or tracked (respectively) by Git.

### Ignoring Files
Go-EarlyBird can ignore any file pattern listed in the `.ge_ignore` and `.gitignore` files. The `--ignorefile` flag can be used to specify a specific path to a file containing ignore patterns.

### Adjusting Severity of A Given Category
Go-Earlybird supports adjusting the severity of a particular category of finding based on patterns that can apply to the filename or the detected match.
An example of when this might be useful could be reducing the severity of the password-secret category when these findings are found in a test directory.
This configuration is done via the `earlybird.json` config file, under the property `adjusted_severity_categories_patterns`. An example of a possible
configuration might be

```json
  "adjusted_severity_categories_patterns": [
    {
      "category": "password-secret",
      "patterns": [
        "(?i)/test/",
        "(?i)/tests/",
        "(?i)/__tests__/",
      ],
      "adjusted_display_severity": "medium",
      "use_filename": true
    }
  ]
```

`adjusted_severity_categories_patterns` is a list of objects with a required `category` field, required `patterns` which are a list of
regular expressions, the required `adjusted_display_severity`, and finally two optional fields `use_filename` and `use_line_value`.
These two fields determine which part of the hit to apply the regular expression patterns. If `use_filename` is true, the match will
be performed on the filename for the given hit. If `use_line_value` is true the match will be performed against the full line value of the hit.
If neither `use_line_value` or `use_filename` are specified, or they are both false, the match will be performed against the exact match of the hit.

### Performing a inclusivity scan

```bash
go-earlybird -path /dir/to/scan -enable inclusivity-rules --display-severity=info
```

### Performing a scan with only certain modules enabled:

```bash
go-earlybird -path /dir/to/scan -enable password-secret -enable content -enable inclusivity-rules
```

## License
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.txt).

## Code of Conduct
This project adheres to the [American Express Community Guidelines](./CODE_OF_CONDUCT.md). By participating, you are expected to honor these guidelines.


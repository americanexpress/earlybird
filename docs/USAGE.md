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
    	Output format [ console | json | csv ] (default "console").
  -git string
    	Full URL to a git repo to scan e.g. github.com/user/repo
  -git-branch string
        Name of branch to be scanned
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
  -ignore-failure
        Avoid the exit code 1 in case of scanner finds valid findings and meets fail threshold
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
  -strict-jks
        Checks for private keys in the JKS file and only return finding if found. If not passed, it will flag jks file. Default is false.
  -suppress
    	Suppress reporting of the secret found (important if output is going to Slack or other logs)
  -update
    	Update module configurations
  -verbose
    	Reports details about file reads
  -version
    	Display version information and exit
  -with-console
        Prints findings in console with JSON format report
  -workers int
    	Set number of workers. (default 100)
  -worksize int
    	Set Line Wrap Length. (default 2500)
  -list-available-modules
    	List available scanning modules. This is useful when inporting configurations from a file and you want to know the module available to configure.
  -module-config-file string
        Absolute path to a json or yaml file for per module level config -- {"modules": { "aModule": { "display_severity": "medium" } } }
  ```

### Performing a scan with only certain modules enabled:

```bash
go-earlybird -path /dir/to/scan -enable password-secret -enable content -enable inclusivity-rules
```
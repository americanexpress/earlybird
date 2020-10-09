# Running EarlyBird

## Standalone
Assuming the `install.sh` script was run, you can kick off the application by running `go-earlybird`, for Windows users call `go-earlybird.exe` instead. Earlybird will scan your current directory unless provided with a target path using the `--path` flag.


## Streamed / Piped input
Using the `--stream` flag, users can stream or pipe file contents to `go-earlybird`.

```
ᐅ go-earlybird -stream < /path/to/file
```
... or:
```
ᐅ cat /path/to/file | go-earlybird --stream
```


## Local Git Scanning
With the flag `--git-staged` or `--git-tracked`, EarlyBird can limit its scan to only look at files that are staged or tracked (respectively) by Git.


## Remote Git Repository Scanning
Using the `--git` flag, you can scan a git repository URL. EarlyBird will clone it into a temporary directory and return results.

For private repositories, please include the `--git-user` flag. You will automatically be prompted for your password. To skip these flags for conveniance or automation, you can use environment variables `gituser` and `gitpassword`

```
ᐅ go-earlybird --git=https://github.com/americanexpress/earlybird
```


## EarlyBird flags
```
Usage of ./go-earlybird:
  -config string
    	Directory where configuration files are stored (default "/Users/jdoe/.go-earlybird/")
  -display-confidence string
    	Lowest confidence level to display [ critical | high | medium | low ] (default "high")
  -display-severity string
    	Lowest severity level to display [ critical | high | medium | low ] (default "medium")
  -enable value
    	Enable individual scanning modules [ ccnumber | common | content | entropy | filename ]
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
    	Full URL to a github organization to scan e.g. github.com/org
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
    	Patterns File (including wildcards) for files to ignore.  (e.g. *.jpg) (default "/Users/jdoe/.ge_ignore")
  -max-file-size int
    	Maximum file size to scan (in bytes) (default 10240000)
  -path string
    	Directory to scan (defaults to CWD) -- ABSOLUTE PATH ONLY (default "/Users/jdoe/Documents/opensource-earlybird/binaries")
  -show-full-line
    	Display the full line where the pattern match was found (warning: this can be dangerous with minified script files)
  -show-rules-only
    	Display rules that would be run, but do not execute a scan
  -show-solutions
    	Display recommended solution for each finding
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
```

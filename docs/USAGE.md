## <a name="running"></a> Running Go-EarlyBird

### <a name="standalone"></a> Standalone
Assuming the setup script was run, you can kick off the application by running 'go-earlybird' / 'go-earlybird.exe' / 'go-earlybird-linux (mac / windows / linux).  See the *Usage* section below

If Go is installed, the project can be downloaded and run with `go run go-earlybird.go`

### Streamed / Piped input
Using the `-stream` flag, users can stream or pipe file contents to 'go-earlybird'.

<details>
<summary>macOS/Linux</summary>

```shell
ᐅ go-earlybird -stream < /path/to/file
ᐅ cat /path/to/file | go-earlybird -stream
```

</details>

<details>
<summary>Windows</summary>

```shell
> go-earlybird -stream > \path\to\file
> type \path\to\file | go-earlybird -stream
```

</details>

### HTTP API

```shell
go-earlybird --http 0.0.0.0:3000
```


`/scan` will accept a multi-part upload and scan the contents, returning json output.

The normal HTTP listener will operate on HTTP/1.1. Go-EarlyBird can be run as HTTPS/2 with the `-https [ip:port]` flag. Note that this also requires the `-https-cert [/path/to/cert]` and `-https-key [/path/to/key]` parameters.

The simple webserver configuration file can be found in the local config directory:
- macOS/Linux: `~/.go-earlybird/webserver.json`
- Windows: `C:\Users\[me]\AppData\go-earlybird\webserver.json`

A separate config file can be specified using the `-http-config [/path/to/configfile]` flag.

### Local Git Scanning
With the flag `-git-staged` or `-git-tracked`, Go-EarlyBird can limit its scan to only look at files that are staged or tracked (respectively) by Git.

## Usage(CLI)
The executable can be called from the command line with the following syntax:

```shell
go-earlybird --help
```

Please note that you can download the earlybird repository using the command
```
git clone https://github.com/americanexpress/earlybird.git
cd earlybird
```

# General Options

- **`--path string`**<br>
`Type: String`, `Default: '/home/earlybird'` <br>
Specify the **directory** to scan (absolute path only).

<details>
<summary>macOS/Linux</summary>

```shell
ᐅ go run go-earlybird.go --path=/scanning_dir --config=./config
```

</details>

<details>
<summary>Windows</summary>

```shell
> go run go-earlybird.go --path=.\scanning_dir --config=.\config
```

</details>

# Selective scanning

- **`--enable value`**<br>
`Type: String`, `Default:''` <br>
Enable individual scanning modules. The supported module names are: **`inclusivity-rules`**, **`password-secret`**

- **`--ignore-fp-rules`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled it **ignore false positive** post-process rules.

# Scan Result Display Options

- **`--display-confidence`**<br>
`Type: String`, `Default:low` <br>
The secret scan results are displayed based on a **minimum confidence threshold**. Only findings with a confidence level equal to or higher than the configured threshold are shown.  
The supported confidence levels are: **low**, **medium**, **high** and **critical**.

The **confidence levels** are prioritized in the following order:
```
critical: 1
high:     2
medium:   3
low:      4
```

- **`--display-severity string`**<br>
`Type: String`, `Default:low` <br>
Secret scan results are displayed based on the **minimum severity threshold**. Only findings with a severity equal to or higher than the configured level are shown. 
The supported severity levels are: **info**, **low**, **medium**, **high** and **critical**.

The **severity levels** are prioritized in the following order:
```
critical: 1
high:     2
medium:   3
low:      4
info:     5
```


- **`--show-full-line`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled it display the full line where the pattern match was found.

- **`--show-rules-only`**<br>
`Type: Boolean`, `Default: false` <br>
Display rules that would be run, but do not execute a scan.


```shell
go run go-earlybird.go  --config=./config --show-rules-only=true
```

- **`--show-solutions`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled true, the scan displays recommended solutions for each finding.

## Scanning considered as failure

- **`--fail-confidence`**<br>
`Type: String`, `Default: low` <br>
Specifies the **minimum confidence level** at which the scan will be considered a failure. If the scanner detects one or more valid findings that meet or exceed this threshold, the scan fails and the process exits with status code 1.
- **`--fail-severity`**<br>
`Type: String`, `Default: low` <br>
Specifies the **minimum severity level** at which the scan will be considered a failure. If the scanner detects one or more valid findings that meet or exceed this threshold, the scan fails and the process exits with status code 1.

# File and Output Options

- **`--file`**<br>
Specify **output file** to store the scan results. 

<details>
<summary>macOS/Linux</summary>

```shell
ᐅ go run go-earlybird.go --path=/scanning_dir --config=./config --file=/myfile.csv
```

</details>

<details>
<summary>Windows</summary>

```shell
> go run go-earlybird.go --path=.\scanning_dir --config=.\config --file=.\myfile.csv
```

</details>

- **`--format`**<br>
`Type: String`, `Default: console` <br>
Specifies the **output format** for the scan results. 
The supported severity levels are: **`console`**, **`json`**, **`csv`**

# Git Repository Scanning

We can do the earlybird scan for the **git repository** also instead of scanning the directory. This could be done by passing the git repo url and respective branch name instead of passing the scanning directory path.

- **`--git`**<br>
`Type: String` <br>
Specifies the **full URL** of the Git repository to scan.

- **`--git-branch string`**<br>
`Type: String` <br>
Specifies name of the **branch** to be scanned.

<details>
<summary>macOS/Linux</summary>

```shell
ᐅ export gituser=<git_user_name>
ᐅ export gitpassword=<git_pat_token>
ᐅ go run go-earlybird.go --config=./config --git=https://github.com/user/repo --git-branch=branch_name
```

</details>

<details>
<summary>Windows</summary>

```shell
> set gituser=<git_user_name>
> set gitpassword=<git_pat_token>
> go run go-earlybird.go --config=.\config --git=https://github.com/user/repo --git-branch=branch_name
```

</details>

- **`--git-commit-stream`**<br>
`Type: Boolean`, `Default: false` <br>
When enabled use **stream IO of Git commit log** as input instead of files. It allows the Earlybird scan to process data directly from a stream (e.g., standard input) instead of scanning files in a directory. 
Specifies name of the branch to be scanned.


```shell
git log |  go run go-earlybird.go --config=./config --git-commit-stream=true 
```


- **`--git-project string`**<br>
`Type: String` <br>
Specify the full URL to a **GitHub organization** or **Bitbucket project** to scan instead of scanning just a repo. Please note that you need to provide the git user-name and the password also during the scanning.

```shell
go run go-earlybird.go --config=./config --git-project=https://github.com/org  --git-user=user_name
```


- **`--git-staged`**<br>
`Type: Boolean` `Default: false`<br>
When enabled, it Scans only **Git staged files** instead of scanning the whole file.

- **`--git-tracked`**<br>
`Type: Boolean` `Default: false`<br>
When enabled, it scan only **Git tracked files** instead of scanning the whole file.

- **`--git-user`**<br>
`Type: String` <br>
Specify the authorized **username** for private Git repositories scanning.

# Earlybird as an HTTP/HTTPS Server
Instead of directly scanning the repository or directory we can **create an HTTP/HTTPS server** also and send multiple requests to scan the repositories.

- **`--disable-keep-alives`**:<br>
`Type: Boolean`, `Default: false` <br>
It is a configuration utility for the earlybird http-Server that controls whether HTTP keep-alives are enabled. By default, keep-alives is enabled, but they can be disabled to close idle connections and conserve resources. This is particularly useful in environments with limited resources.

- **`--http`**<br>
`Type: String` <br>
Specify an **IP and Port** while creating an HTTP server to listen to the API (e.g., `127.0.0.1:8080`).
- **`--http-config`**<br>
`Type: String` <br>
Specify the **path to webserver config** JSON file. (e.g., `my-file.json`). 
```shell
go run go-earlybird.go --http=127.0.0.1:8080 --http-config=my-file.json
```
The default json configuration looks as listed below:
```shell
{
   WriteTimeout: 60,
   ReadTimeout:  60,
   IdleTimeout:  120,
}
```

- **`--https`**<br>
`Type: String` <br>
Specify an **IP and Port** while create an HTTPS/2 server to listen to the API (e.g., `127.0.0.1:8080`). Note, please also add the https-cert and https-key using the cli flags listed below.

- **`--https-cert`**<br>
`Type: String` <br>
Specify the Certificate file for TLS connect for the https server.
- **`--https-key`**<br>
`Type: String` <br>
Specify the private key file for TLS connect for the https server.

<details>
<summary>macOS/Linux</summary>

```shell
ᐅ go run go-earlybird.go --https=127.0.0.1:8080 --https-cert=/path/to/cert --https-key=/path/to/key
```

</details>

<details>
<summary>Windows</summary>

```shell
> go run go-earlybird.go --https=127.0.0.1:8080 --https-cert=\path\to\cert --https-key=\path\to\key
```

</details>

# Additional Options

- **`--ignorefile`**<br>
`Type: String`, `Default: './.ge_ignore` <br>
Specify the patterns file for files to ignore. You can list files inside the ge_ignore as listed below:
```shell
/folder-name1/**
/folder-name2/subfolder/**
**/file-name.extension
```


- **`--ignore-failure`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled it avoid **exit code 1** if valid findings meet fail threshold.

- **`--max-file-size`**<br>
`Type: Int`, `Default: 10240000` <br>
Specify the maximum file size in bytes to scan.
- **`--module-config-file`**<br>
`Type: String`<br>
Specify the path to file with per-module config settings.

- **`--skip-comments`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled it Skips scanning comments in files (applies only to the `content` module).
- **`--stream`**<br>
`Type: Boolean`, `Default: false` <br>
When enabled true the Earlybird scan uses stream IO as input instead of files. So instead of passing the scanning directory we passes a stream data source.


```shell
echo 'my file has  access_key = secret_key' |  go run go-earlybird.go --config=./config --stream=true
```


- **`--suppress`**<br>
`Type: Boolean`, `Default: false` <br>
When enabled true it mask the secret line value(*****) while reporting of the secret found.
- **`--update`**<br>
`Type: Boolean`, `Default: false` <br>
When set to true it updates module configurations.
- **`--verbose`**<br>
`Type: Boolean`, `Default: false` <br>
When set to true the scan report comes with details about file reads.
-  **`--version`**<br>
`Type: Boolean`, `Default: false` <br>
Display version information and then exit.
- **`--with-console`**<br>
`Type: Boolean`, `Default: false` <br>
If enabled it allows using along with the `--format` flag, to print findings report in console with JSON format.
- **`--workers int`**<br>
`Type: Int`, `Default: 100` <br>
Set number of workers needed to create the worker pool while scanning the code.
- **`--worksize`**<br>
`Type: Int`, `Default: 2500` <br>
Set line wrap length. When a lengthy line is fed as an input from the code for the scan work-size define the max character size of line for scanning, remaining character get fed as a new input with the same process, Please note that there is also an overlap happens to not miss any secrets that got split.
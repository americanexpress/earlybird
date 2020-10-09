# EarlyBird REST API
The normal HTTP listener will operate on HTTP/1.1, which can be run with the `http [ip:port]` flag.  EarlyBird can be run as HTTPS/2 with the `--https [ip:port]` flag.  Note that this also requires the `--https-cert [/path/to/cert]` and `--https-key [/path/to/key]` parameters.

Here's an example on how to start the API from the command line.

```
ᐅ go-earlybird --http 0.0.0.0:3000
```

## REST API Endpoints
`/scan` will accept a multi-part upload and scan the contents, returning JSON output.
`/scan/git?url=https://example.com/repo.git` will accept a git repository URL, clone and scan the contents, returning JSON output.
`/labels` will return all of the labels from the config files as a JSON output
`/categories` will return all of the rule categories from the config files as a JSON output
`/categorylabels` will return all of the labels per category from the config files as a JSON output


The simple webserver configuration file can be found in the local config directory (`~/.go-earlybird/webserver.json` or `C:\Users\[user]\AppData\go-earlybird\webserver.json`).  A separate config file can be specified using the `--http-config [/path/to/configfile]` flag.

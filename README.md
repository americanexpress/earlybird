![Logo](docs/GoEarlyBird-logo_sm.png)

EarlyBird is a sensitive data detection tool capable of scanning source code repositories for clear text password violations, PII, outdated cryptography methods, key files and more. It can be used to scan remote git repositories, local files or directories or as a pre-commit step.

<br />

## Installation

### Docker
**Running `docker build -t earlybird .` will automate the following steps**
- Pull `golang:1.18-alpine` image used to build EarlyBird binary.
- Install packages neccissary to compile EarlyBird binary (`gcc`/`libc-dev`)
- Update `mod.go` and download required dependencies (`go mod tidy`)
- Run unit tests found in `./pkg/...` aborting with exit code `1 ` if a failure is detected.
- Generate a `Linux` binary in `./binaries/`
- Pull `alpine:latest` image used to run EarlyBird
- Copy `./config/*` to configuration directory `/root/.go-earlybird`
- Copy `.ge_ignore` to configuration directory `/root/.go-earlybird`
- Copy `./binaries/go-earlybird` binary to install directory

> Note: If environment variable `EBVERSION`is defined, it's used as version tag of the binary, otherwise `dev` will be used (`buildflags.Version=%EBVERSION%`)

<br />

**Build image and run scan against local filesystem**

```bash
docker build -t earlybird .
docker run --rm -it -v /local/filesystem/path/:/app/ earlybird 
```
**Make full use of command line arguments by changing the entrypoint and including custom configs**

```bash
docker run --rm -it -v /local/filesystem/path:/app/ -v /path/to/local/config:/root/.go-earlybird/config --entrypoint go-earlybird earlybird "-display-severity=high -fail-severity=high"
```

> Note: THe EarlyBird container will run against its local `/app/` directory. Mount filesystems to the `/app/` directory to make use of this.

> Note: To customize the default arguments within the EarlyBird docker image, alter the `CMD` line at the end of the `Dockerfile` with desired arguments.


<br />



### Linux & Mac
**Running `build.sh` will automate the following steps**
- Update `mod.go` and download required dependencies (`go mod tidy`)
- Run unit tests found in `./pkg/...` aborting with exit code `1 ` if a failure is detected.
- Generate a `Linux` binary in `./binaries/`
- Generate a `MacOS` binary in `./binaries/`
- Generate a `Windows` binary in `./binaries/`
> Note: If an argument is specified, it is used as the version tag, otherwise `dev` will be used (`buildflags.Version=${EBVERSION}`)

<br />

**Running `install.sh` will automate the following steps**
- Create install directory if it doesn't already exist (`${HOME}\.go-earlybird`)
- Copy `./config/*` to install directory
- Copy `.ge_ignore` to install directory
- Copy `go-earlybird` binary to `/usr/local/bin/`
- Ensure `go-earlybird` binary in `/usr/local/bin/` is executable

<br />

**Build binaries and install as current user**
> Note: User must have write access to /usr/local/bin/
```
./build.sh && ./install.sh
```

<br />

### Windows
**Running `build.bat` will automate the following steps**
- Update `mod.go` and download required dependencies (`go mod tidy`)
- Run unit tests found in `./pkg/...` aborting with exit code `1 ` if a failure is detected.
- Generate a `Linux` binary in `./binaries/`
- Generate a `MacOS` binary in `./binaries/`
- Generate a `Windows` binary in `./binaries/`
> Note: If an argument is specified, it is used as the version tag, otherwise `dev` will be used (`buildflags.Version=%EBVERSION%`)

<br />

**Running `install.bat` will automate the following steps**
- Create install directory if it doesn't already exist (`%APPDATA%\go-earlybird`)
- Copy `.\config\*` to install directory
- Copy `.ge_ignore` to install directory
- Copy `go-earlybird` binary to install directory
- Add install directory to user `PATH` variable (`HKEY_CURRENT_USER\Environment`)

<br />


**Build binaries and install as current user**
```
.\build.bat && .\install.bat
```

<br />

## Uninstallation

### Docker
**Uninstall EarlyBird image**
```
docker image rm earlybird
```

<br />


### Linux & Mac
**Running `uninstall.sh` will automate the following steps**
- Remove build logs potentially left over from `build.bat` (`${HOME}/eb-build-log.log`)
- Remove install directory (`${HOME}/.go-earlybird`)
- Remove binary (`/usr/local/bin/go-earlybird`)

<br />

**Uninstall EarlyBird**
```
./uninstall.sh
```
<br />

### Windows
**Running `uninstall.bat` will automate the following steps**
- Remove build logs potentially left over from `build.bat` (`%APPDATA%\eb-build-log.log`)
- Remove install directory (`%APPDATA%\go-earlybird`)
- Remove install directory from user `PATH` variable (`HKEY_CURRENT_USER\Environment`)

<br />

**Uninstall EarlyBird**
```
.\uninstall.bat
```

<br />

## Usage
Launch a basic EarlyBird scan against a directory:
```bash
# Linux
go-earlybird --path=/path/to/directory

# Windows
go-earlybird.exe --path=C:\path\to\directory
```
Launch a basic EarlyBird scan against a file:
```bash
# Linux
go-earlybird --path=/path/to/file.txt

# Windows
go-earlybird.exe --path=C:\path\to\file.txt
```

Scan a remote git repo:
```
go-earlybird --git=https://github.com/americanexpress/earlybird
```
[Click here for Detailed Usage instructions.](./docs/USAGE.md)

<br />

## Documentation
 - [Usage - How do I use Earlybird?](./docs/USAGE.md)
 - [Modules - What is a Module? How do I create one?](./docs/MODULES.md)
 - [Hooks - How do I use Earlybird as Pre-Commit Hook?](./docs/HOOKS.md)
 - [REST API - How do I use Earlybird as REST API?](./docs/REST.md)
 - [False Positives - How are they managed? How do I filter them?](./docs/FALSEPOSITIVES.md)
 - [Labels - What are labels? How do I create my own?](./docs/LABELS.md)
 - [Ignore - How do I skip lines or files intentionally?](./docs/IGNORE.md)
 - [Inclusivity - How do I perform an inclusivity scan?](./docs/INCLUSIVITY.md)

<br />

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

<br />

## Contributing
We welcome your interest in the American Express Open Source Community on Github. Any Contributor to
any Open Source Project managed by the American Express Open Source Community must accept and sign
an Agreement indicating agreement to the terms below. Except for the rights granted in this 
Agreement to American Express and to recipients of software distributed by American Express, You
reserve all right, title, and interest, if any, in and to your contributions. Please
[fill out the Agreement](https://cla-assistant.io/americanexpress/earlybird).

<br />

## License
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.txt).

<br />

## Code of Conduct
This project adheres to the [American Express Community Guidelines](./CODE_OF_CONDUCT.md). By participating, you are expected to honor these guidelines.


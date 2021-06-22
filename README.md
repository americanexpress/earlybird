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
 - [Inclusivity - How do I perform an inclusivity scan?](./docs/INCLUSIVITY.md)


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

## License
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.txt).

## Code of Conduct
This project adheres to the [American Express Community Guidelines](./CODE_OF_CONDUCT.md). By participating, you are expected to honor these guidelines.


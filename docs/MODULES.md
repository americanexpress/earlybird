# EarlyBird Modules
Modules are configurable rule sets that define patterns and target areas for searching within source code files. The currently implemented modules include:

## Included Modules:
 - __File Names (filename)__: Scan the file list recursively, looking for filename patterns that would indicate credentials, keys, and sensitive PII.  We're looking for things like `id_rsa`, things that end in `pem`, etc.
 - __File Content Patterns (content)__: Looks for patterns within the contents of files, things like `password: `, and `BEGIN RSA PRIVATE KEY` will pop up here.  Other types of sensitive PII data elements and secrets will be detected as well, such as IBAN, SSN, IP Addresses, Email Addresses, Phone Numbers, etc.  This also looks for insecure cryptographic algorithms and pseudo-random number generation, as well as suspicious comments like "HACK" and "FIXME".
 - __File Content Entropy (entropy)__:  Scan files for strings with high (Shannon) entropy, which could indicate passwords or secrets stored in the files, for example: `kwaKM@£rFKAM3(a2klma2d`
 - __Credit Card Numbers (ccnumber)__:  Scan files for strings that match major credit card number patterns.  Any potential hits are passed through a Luhn/mod10 check to verify that they are valid card numbers, and all numbers that are identified as designated test values are ignored.
 - __Commonly Used / Default Passwords (common)__:  Scan files for default and commonly used/abused passwords.
 &nbsp;
 
## Creating New Modules:
New modules can be added via json rules files in the user's `go-earlybird` configuration directory.  Simply add a new json file into this directory (e.g. `custom-rules.json`) with the following structure, and EarlyBird will detect and load the rules.  Keeping these custom rules in a separate file will ensure they do not get overwritten when EarlyBird is updated.
```
{
    "Searcharea": "<Where in the file should the scan search?  Supports `body` or `filename`>",
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

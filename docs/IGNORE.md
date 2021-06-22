# Ignoring and Skipping in EarlyBird
EarlyBird scans all files in the target directory or uploaded via the API. We understand you may want to skip some unessential files such as images. We built in functionality for an ignore file or an ignore comment. Here's how to use them below.


### Ignoring Files
EarlyBird can ignore any file pattern listed in the `.ge_ignore` and `.gitignore` files. The `--ignorefile` flag can be used to specify a specific path to a file containing ignore patterns.


### Ignoring Lines
Annotations can be used in any file through comments or any other text value to flag the line to be ignored.  If a file will intentionally contain a potential secret (e.g. test data), you can specify `EARLYBIRD-IGNORE` in the line and the scan will skip it.  See the example below:

```
public String get_test_pass() // return test case password
{
    String test_pass = "unit_test"; //EARLYBIRD-IGNORE
    return test_pass;
}
``` 

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

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

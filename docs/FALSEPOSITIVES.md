# False Positive Detection
EarlyBird has a rules engine for excluding false positives from the results.  Each rule in `false-positives.json` is tied to one or more scan rules (using the `Codes` field).

The `Pattern` field is a regular expression that is evaluated against any hit that matches the `Code`, as long as the file containing that hit has an extension matching a value in the `FileExtensions` value (if that value is empty, all file extensions will be considered).

```
{
     "rules": [
         {
             "Codes": [3013],
             "Pattern": "(000-000-0000)",
             "FileExtensions": [],
             "Description": "Ignore a false positive phone number"
         },
         {
             "Codes": [3058,3069,3070,3071,3072,3073,3074,3075],
             "Pattern": ".*",
             "FileExtensions": [".md", ".txt", ".doc", ".pdf", ".docx", ".csv", ".html", ".htm"],
             "Description": "Ignore deprecated crypto in documents"
         },
         ...
```

### In the examples above:
1. Any hit found with rule 3013 (looking for 10-digit phone number patterns), that matches all zeroes (000-000-0000), in any type of file will be ignored
2. Any hit found with in the specified rules (looking for indicators of deprecated crypto method use like 3DES or MD5) in document files will be ignored

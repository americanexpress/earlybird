# EarlyBird as a git pre-commit hook
EarlyBird can easily be added to a bash pre-commit hook script, as seen below *(Remember -- failures in pre-commit hooks serve as a warning.  If you need to commit the code despite this warning, you can override it with the `--no-verify` flag, but it is highly recommended that _all_ secrets stay out of git repositories)*:
Add the following file to ```./git/hooks/pre-commit``` Please ensure your pre-commit hook has executable permissions ```chmod +x ./git/hooks/pre-commit```
```bash
#!/bin/sh

echo "Running Go-EarlyBird pre-commit hook"
go-earlybird -display-severity=high -fail-severity=high -git-staged

# $? stores exit value of the last command
if [ $? -ne 0 ]; then
 echo ""
 echo "Secrets detection tests must pass before commit!"
 echo ""
 exit 1
fi
```
> NOTE To run the pre-commit hook without failing the commit, you can remove the `exit 1` line, although we recommend keeping the commit-blocking in place)

<br />

The following is an example of running EarlyBird from package.json against a local repository before commit, failing the commit if high or critical issues are found (requires git-bash on windows):
```json
{
  "name": "demoproj",
  "version": "1.0.0",
  "description": "demo",
  "main": "app.js",
  "dependencies": {
    "express": "4.16.2"
  },
  "devDependencies": {
    "pre-commit": "^1.2.2"
  },
  "scripts": {
    "earlybird:pre-commit": "go-earlybird --fail-severity=high --git-staged --format=json --file=pre-commit-output.log"
  },
  "pre-commit": [
    "earlybird:pre-commit"
  ],
  "author": "",
  "license": "ISC"
}
```

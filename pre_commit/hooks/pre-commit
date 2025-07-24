#!/usr/bin/env bash
# Using unofficial bash strict mode
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# Check if earlybird is installed, skipping the check if it is not installed
if [[ ! -x $(command -v go-earlybird) ]]; then
  echo "earlybird is not installed. Skipping pre-commit check..."
  echo "To install, follow the guide here:"
  echo "https://github.com/americanexpress/earlybird?tab=readme-ov-file#installation"
  exit 0
fi

ARGS=${*:-}
echo "Running Go-EarlyBird pre-commit hook"
go-earlybird -git-staged ${ARGS}

if [ $? -ne 0 ]; then
 echo "Earlybird scan detected possible secrets in your staged changes."
 echo "If you think this might be a false positive, you can exclude secret values by following this guide:"
 echo "https://github.com/americanexpress/earlybird/blob/40a7a9b8b43d0fe6b510359ab6cff0675dee468c/docs/FALSEPOSITIVES.md"
 exit 1
fi
# [4.1.0](https://github.com/americanexpress/earlybird/compare/v4.0.0...v4.1.0) (2024-02-07)


### Features

* **config:** added config to avoid exit code 1 for scan failure ([5b0e14a](https://github.com/americanexpress/earlybird/commit/5b0e14a039bf5d2eeebc71867864a96f2d9dd3d0))

# [4.0.0](https://github.com/americanexpress/earlybird/compare/v3.16.0...v4.0.0) (2023-09-19)


### Bug Fixes

* **nameScanner:** updated logic to fail a scan to respect severity and confidence ([a7b9a56](https://github.com/americanexpress/earlybird/commit/a7b9a5684eab2b429dc9698e48cbf1816bf82fe4))


* Merge pull request #82 from americanexpress/feat/version-upgrade ([810b7fb](https://github.com/americanexpress/earlybird/commit/810b7fb6c0c66a8ea77deb05df4777f732d8ac6e)), closes [#82](https://github.com/americanexpress/earlybird/issues/82)


### BREAKING CHANGES

* Deps update for Vulnerability fixes, version update to go 1.20 and module update

# [3.16.0](https://github.com/americanexpress/earlybird/compare/v3.15.0...v3.16.0) (2023-07-10)


### Features

* **config:** not exiting process after update and reloading config ([#74](https://github.com/americanexpress/earlybird/issues/74)) ([a96d2a2](https://github.com/americanexpress/earlybird/commit/a96d2a2567ab3cf186e3309da20031dc749e241f))

# [3.15.0](https://github.com/americanexpress/earlybird/compare/v3.14.0...v3.15.0) (2023-06-01)


### Features

* **file-with-console:** updated code based to broadcast and write it… ([#60](https://github.com/americanexpress/earlybird/issues/60)) ([92e6123](https://github.com/americanexpress/earlybird/commit/92e6123fec5a7d6046f8c5914796550050522df5))

# [3.14.0](https://github.com/americanexpress/earlybird/compare/v3.13.2...v3.14.0) (2023-03-16)


### Features

* add arm64 binary for use on M1 hardware ([23c4ff0](https://github.com/americanexpress/earlybird/commit/23c4ff0b79180f8cff141c7e5c38101de5958e31))

## [3.13.2](https://github.com/americanexpress/earlybird/compare/v3.13.1...v3.13.2) (2023-02-23)


### Bug Fixes

* **config:** updated config based display and failure ([#55](https://github.com/americanexpress/earlybird/issues/55)) ([5247ccf](https://github.com/americanexpress/earlybird/commit/5247ccf4f438f9b3086ff7c16e70f2e46d0ff9a6))

## [3.13.1](https://github.com/americanexpress/earlybird/compare/v3.13.0...v3.13.1) (2022-10-03)


### Bug Fixes

* scan failure provision for info level hits ([1174e45](https://github.com/americanexpress/earlybird/commit/1174e45400dd375d0f6555e042d8ce7959b52e61))

# [3.13.0](https://github.com/americanexpress/earlybird/compare/v3.12.0...v3.13.0) (2022-09-07)


### Bug Fixes

* release branch set to main ([fcdba6f](https://github.com/americanexpress/earlybird/commit/fcdba6f995c3e699e21cffe8fa33d132771c7c70))
* remove the default ignore with .gitignore to avoid missing scans on the repo. ([2e6255d](https://github.com/americanexpress/earlybird/commit/2e6255d0aaa79821902ede4e90a41e2e10cdd4d4))


### Features

* add keepAlive flag and fix worker flag read ([#47](https://github.com/americanexpress/earlybird/issues/47)) ([549081d](https://github.com/americanexpress/earlybird/commit/549081d257a0d2de4a9f256e1d9a948d2a670c30))
* add ldflags for version injection during artifacts build ([32fc145](https://github.com/americanexpress/earlybird/commit/32fc14532597334c6b99900d4b092cd100768632))
* use semantic versioning to create semantic releases and changelog.md file ([2e8a2e9](https://github.com/americanexpress/earlybird/commit/2e8a2e91cf0f1f8ccd4b96d01c9a5f5db0c06cd8))

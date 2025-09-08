## [4.4.1](https://github.com/americanexpress/earlybird/compare/v4.4.0...v4.4.1) (2025-09-08)


### Bug Fixes

* basic auth fp ([30071ec](https://github.com/americanexpress/earlybird/commit/30071ec2c0b34e96ac81a0803680dba0ce29031b)), closes [#178](https://github.com/americanexpress/earlybird/issues/178)
* basic auth fp ([a8a6e97](https://github.com/americanexpress/earlybird/commit/a8a6e97ffa5a2eb663a94e4add6a4d30b1049e0e))
* eb release ([ac83277](https://github.com/americanexpress/earlybird/commit/ac832778cbecffc11a0618a7473753db36d4553e))

# [4.4.0](https://github.com/americanexpress/earlybird/compare/v4.3.4...v4.4.0) (2025-09-03)


### Features

* adding postprocess for basic auth header ([#163](https://github.com/americanexpress/earlybird/issues/163)) ([b803750](https://github.com/americanexpress/earlybird/commit/b803750f62e81b05ad2c59316dc70f6461738f80))

## [4.3.4](https://github.com/americanexpress/earlybird/compare/v4.3.3...v4.3.4) (2025-08-28)


### Bug Fixes

* making the addlistener channel as buffered ([986daa7](https://github.com/americanexpress/earlybird/commit/986daa728aacb2468cca7361cc616d9058ad4c8b))
* making the addlistener channel as buffered ([3e987dd](https://github.com/americanexpress/earlybird/commit/3e987dd8ed8b4cea1e7ac97bc8d6ba53790a06ea))

## [4.3.3](https://github.com/americanexpress/earlybird/compare/v4.3.2...v4.3.3) (2025-08-12)


### Bug Fixes

* addressing the deadlock issue ([#150](https://github.com/americanexpress/earlybird/issues/150)) ([c0e61c6](https://github.com/americanexpress/earlybird/commit/c0e61c6c81881ee15e28f0a34f0fb044b82d0471))

## [4.3.2](https://github.com/americanexpress/earlybird/compare/v4.3.1...v4.3.2) (2025-08-08)


### Bug Fixes

* reuse bytes.Buffer in loop to reduce allocations ([#154](https://github.com/americanexpress/earlybird/issues/154)) ([8693c7f](https://github.com/americanexpress/earlybird/commit/8693c7f701b1f8d17753b423a34a236d1f8d15ef))

## [4.3.1](https://github.com/americanexpress/earlybird/compare/v4.3.0...v4.3.1) (2025-07-22)


### Bug Fixes

* simplify if statement in skipAccountTokenPass ([#149](https://github.com/americanexpress/earlybird/issues/149)) ([347c4ba](https://github.com/americanexpress/earlybird/commit/347c4ba69c3d3472500c538d8568b111f37a68c1))

# [4.3.0](https://github.com/americanexpress/earlybird/compare/v4.2.6...v4.3.0) (2025-07-10)


### Features

* adding secret detection in .pem file ([#143](https://github.com/americanexpress/earlybird/issues/143)) ([a496681](https://github.com/americanexpress/earlybird/commit/a49668182066ac943cfea2f6030311945ea5b540))

## [4.2.6](https://github.com/americanexpress/earlybird/compare/v4.2.5...v4.2.6) (2025-07-10)


### Bug Fixes

* remove CC testcase ([6778208](https://github.com/americanexpress/earlybird/commit/6778208df0626bb7dc2c83ec754312ae11d4bf97))
* Updating the file extension matching ([6f2f4bf](https://github.com/americanexpress/earlybird/commit/6f2f4bf4c8f3684392a303092efe89e2e4c7a866))

## [4.2.5](https://github.com/americanexpress/earlybird/compare/v4.2.4...v4.2.5) (2025-05-19)


### Bug Fixes

* log error and continue when error reading compressed file ([7535d3a](https://github.com/americanexpress/earlybird/commit/7535d3aeea88fe2361107d0f3833a49ae5dca718))

## [4.2.4](https://github.com/americanexpress/earlybird/compare/v4.2.3...v4.2.4) (2025-04-23)


### Bug Fixes

* **module:** updated code to report finding based on module config ([4f2561c](https://github.com/americanexpress/earlybird/commit/4f2561cd64b0e8b42c5a4a3614746ebba9bff070))

## [4.2.3](https://github.com/americanexpress/earlybird/compare/v4.2.2...v4.2.3) (2025-03-19)


### Bug Fixes

* detect unicode char in post processor ([02c8307](https://github.com/americanexpress/earlybird/commit/02c8307ce4c400f79ed3d25867b9e18b22cf737f))

## [4.2.2](https://github.com/americanexpress/earlybird/compare/v4.2.1...v4.2.2) (2025-02-17)


### Bug Fixes

* **jks-http:** updated code to check for pk in jks via http ([f12afaa](https://github.com/americanexpress/earlybird/commit/f12afaa29d8e002f61241469869f087e8ead8e3a))

## [4.2.1](https://github.com/americanexpress/earlybird/compare/v4.2.0...v4.2.1) (2025-02-12)


### Bug Fixes

* **ssn:** updated ssn rule to ignore flaging 9xx in first group ([a6ea533](https://github.com/americanexpress/earlybird/commit/a6ea533c120ed3d3457d6ad71cc27ee573df1821))

# [4.2.0](https://github.com/americanexpress/earlybird/compare/v4.1.3...v4.2.0) (2025-02-04)


### Bug Fixes

* **option:** removed options from parsing since not required ([013b439](https://github.com/americanexpress/earlybird/commit/013b439265a367317b7dc6dd8ff2f316addd620f))


### Features

* **jks-check:** added config to strictly check for key in jks ([7ce4e75](https://github.com/americanexpress/earlybird/commit/7ce4e75087ae569d0f2a33fca9032c5585bf56fb))

## [4.1.3](https://github.com/americanexpress/earlybird/compare/v4.1.2...v4.1.3) (2025-01-07)


### Bug Fixes

* issue with git tracked file not ignoreing ge_ignore patterns ([d651d7d](https://github.com/americanexpress/earlybird/commit/d651d7dfafc0dd09f276a4d64b8de592327b9aed))

## [4.1.2](https://github.com/americanexpress/earlybird/compare/v4.1.1...v4.1.2) (2024-05-20)


### Bug Fixes

* require digits for card mod10 check ([0462806](https://github.com/americanexpress/earlybird/commit/04628068e642a3c23754e8f3085469968870e09c))
* **ssn:** updated ssn rule to ignore flaging 9xx in third group ([2230817](https://github.com/americanexpress/earlybird/commit/2230817baa0ef6de3918219d9b789b17d4228346))

## [4.1.1](https://github.com/americanexpress/earlybird/compare/v4.1.0...v4.1.1) (2024-02-12)


### Bug Fixes

* the extension regex was not reading the dot ([53cf2c8](https://github.com/americanexpress/earlybird/commit/53cf2c8b28c75160494784f12d574df24f996159))

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

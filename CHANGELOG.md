# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2023-02-14

### Added

- Ability to specify `credential_process` option in the config to source Vault token from an external process

## [0.2.1] - 2020-07-28

### Changed

- Testing is done now by using Vault utilities instead of mocking

### Fixed

- Handling scenario of Vault returning a `nil` secret

## [0.2.0] - 2020-07-22

### Added

- Allow insecure SSL communication with Vault

### Fixed

- Fix wrong log level for exclusion paths

## [0.1.0] - 2020-07-22

First release

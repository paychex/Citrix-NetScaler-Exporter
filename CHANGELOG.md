# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [2.0.1] - 2018-04-30
### Changed
 - Added support for the "query" endpoint to watch multiple Netscaler devices.  Moved the collector code into its own subdirectory to make the code more modular.  Updated README.md to reflect changes in support for multi-query.

## [2.0.0] - 2017-10-10
### Changed
 - Log entries are no longer sent to a file.  Instead they are logged to stdout in logfmt format.

## [1.4.1] - 2017-08-22
### Fixed
 - NetScaler API bug meant that trying to retrieve stats from a service group member which used a wildcard port (65535 in API and CLI, * in GUI) resulted in error.  Skipping these members until the bug is resolved.

## [1.4.0] - 2017-07-25
### Changed
 - Authentication to the NetScaler now only happens once per scrape; the session cookie is saved and re-used in future requests.  When the scrape finishes, the session is disconnected.  Previously each API request was authenticated individually.

## [1.3.0] - 2017-07-21
### Added
 - Exporting Service Group metrics

## [1.2.0] - 2017-07-21
### Added
 - Added service state; 1 if service is up and 0 for any other state.

## [1.1.0] - 2017-07-20
### Added
 - Added model_id which represents, on the VPX line at least, the maximum licensed throughout of the appliance.  For example a VPX 1000 allows for 1000 MB/s.

## [1.0.0] - 2017-07-15
Initial release

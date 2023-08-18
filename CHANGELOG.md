## Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [0.1.0] - 2023-08-18
### Added
- A `Config` type to define the configuration of the Budget microservice.

### Updated
- The `Service` and `NewService` functions to accept a `Config` type.
- The `micro` package to `v0.2.0`.

## [0.0.0] - 2022-04-15
### Added
- Service to make it a mock and testable integration for any microservice 
or API gateway.
- Models all the data structures that can be expected as responses from
the Budget microservice.
- Budget, Group, Item and Event are all exchange functions that format 
and process HTTP exchanges between the Microservice Package (MSP) and the
Budget Microservice.

## [Released]
